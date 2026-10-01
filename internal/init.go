package internal

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nagorneva/nagorneva_reviews/internal/config"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/login"
	filter "github.com/nagorneva/nagorneva_reviews/internal/module/reviews/action/filter_reviews"
	onmap "github.com/nagorneva/nagorneva_reviews/internal/module/reviews/action/get_doctors_on_map"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/auth"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/httpx"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/password"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/storage"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type initializer struct {
	*App
	err    error
	cfg    config.Config
	db     *pgxpool.Pool
	s3     *storage.S3
	token  auth.Service
	crud   crud.Action
	login  login.Action
	filter filter.Action
	onMap  onmap.Action
}

func (a *App) initConfig(_ context.Context) *initializer {
	i := &initializer{App: a}
	i.cfg, i.err = config.Load()
	if i.err == nil {
		i.token = auth.New(i.cfg.JWTSecret, i.cfg.JWTTTL)
	}
	return i
}
func (i *initializer) initPostgres(ctx context.Context) *initializer {
	if i.err != nil {
		return i
	}
	i.db, i.err = pgxpool.New(ctx, i.cfg.DatabaseURL)
	if i.err == nil {
		i.shutdown = i.db.Close
		i.err = i.db.Ping(ctx)
	}
	return i
}
func (i *initializer) initS3(ctx context.Context) *initializer {
	if i.err != nil {
		return i
	}
	i.s3, i.err = storage.New(ctx, i.cfg.S3Endpoint, i.cfg.S3Region, i.cfg.S3Bucket, i.cfg.S3AccessKey, i.cfg.S3SecretKey, i.cfg.S3PublicURL)
	return i
}
func (i *initializer) initModules(ctx context.Context) *initializer {
	if i.err != nil {
		return i
	}
	if err := i.ensureInitialUser(ctx); err != nil {
		i.err = err
		return i
	}
	i.crud = crud.New(i.db)
	i.login = login.New(i.db)
	i.filter = filter.New(i.db)
	i.onMap = onmap.New(i.db)
	return i
}
func (i *initializer) ensureInitialUser(ctx context.Context) error {
	if i.cfg.AdminEmail == "" || i.cfg.AdminPassword == "" {
		return nil
	}
	var exists bool
	if err := i.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, i.cfg.AdminEmail).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil
	}
	hash, err := password.Hash(i.cfg.AdminPassword)
	if err != nil {
		return err
	}
	_, err = i.db.Exec(ctx, `INSERT INTO users(email,password) VALUES($1,$2)`, i.cfg.AdminEmail, hash)
	return err
}
func (i *initializer) initRouter(_ context.Context) *initializer {
	if i.err != nil {
		return i
	}
	r := chi.NewRouter()
	r.Use(middleware.Recoverer, middleware.RequestID, middleware.RealIP, middleware.Timeout(30*time.Second), cors)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { httpx.JSON(w, 200, map[string]string{"status": "ok"}) })
	r.Post("/api/v1/doctors_on_map", i.getDoctorsOnMap)
	r.Post("/api/v1/reviews/filter", i.filterReviews)
	r.Post("/api/v1/admin/login", i.signIn)
	r.Route("/api/v1/admin", func(r chi.Router) {
		r.Use(i.token.Middleware)
		for _, resource := range []string{"users", "cities", "specialties", "doctors", "courses", "reviews"} {
			i.resourceRoutes(r, resource)
		}
		r.Get("/doctors/{id}/reviews", i.doctorReviews)
		r.Post("/doctors/{id}/photo", i.uploadPhoto)
	})
	i.server = &http.Server{Addr: i.cfg.HTTPAddr, Handler: r, ReadHeaderTimeout: 5 * time.Second}
	return i
}
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (i *initializer) resourceRoutes(r chi.Router, resource string) {
	r.Route("/"+resource, func(r chi.Router) {
		r.Get("/", i.list(resource))
		r.Post("/", i.create(resource))
		r.Get("/{id}", i.get(resource))
		r.Patch("/{id}", i.update(resource))
		r.Delete("/{id}", i.delete(resource))
	})
}
func (i *initializer) signIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if httpx.Decode(r, &req) != nil {
		httpx.Error(w, 400, "invalid request")
		return
	}
	id, err := i.login.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		httpx.Error(w, 401, "invalid email or password")
		return
	}
	token, err := i.token.Issue(id)
	if err != nil {
		httpx.Error(w, 500, "token issue failed")
		return
	}
	httpx.JSON(w, 200, map[string]any{"access_token": token, "token_type": "Bearer"})
}
func (i *initializer) list(res string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, e := i.crud.List(r.Context(), res)
		writeAction(w, v, e, 200)
	}
}
func (i *initializer) get(res string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, e := idParam(r)
		if e == nil {
			v, x := i.crud.Get(r.Context(), res, id)
			writeAction(w, v, x, 200)
		} else {
			httpx.Error(w, 400, e.Error())
		}
	}
}
func (i *initializer) create(res string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, e := decodeValues(r)
		if e == nil {
			e = hashPassword(v)
			if e == nil {
				out, x := i.crud.Create(r.Context(), res, v)
				writeAction(w, out, x, 201)
			}
		}
		if e != nil {
			httpx.Error(w, 400, e.Error())
		}
	}
}
func (i *initializer) update(res string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, e := idParam(r)
		if e == nil {
			v, x := decodeValues(r)
			if x == nil {
				x = hashPassword(v)
			}
			if x == nil {
				out, z := i.crud.Update(r.Context(), res, id, v)
				writeAction(w, out, z, 200)
			} else {
				e = x
			}
		}
		if e != nil {
			httpx.Error(w, 400, e.Error())
		}
	}
}
func (i *initializer) delete(res string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, e := idParam(r)
		if e == nil {
			e = i.crud.Delete(r.Context(), res, id)
			if e == nil {
				w.WriteHeader(204)
				return
			}
		}
		writeAction(w, nil, e, 204)
	}
}
func (i *initializer) doctorReviews(w http.ResponseWriter, r *http.Request) {
	id, e := idParam(r)
	if e == nil {
		v, x := i.crud.DoctorReviews(r.Context(), id)
		writeAction(w, v, x, 200)
	} else {
		httpx.Error(w, 400, e.Error())
	}
}
func (i *initializer) uploadPhoto(w http.ResponseWriter, r *http.Request) {
	id, e := idParam(r)
	if e != nil {
		httpx.Error(w, 400, e.Error())
		return
	}
	if i.s3 == nil {
		httpx.Error(w, 503, "S3 is not configured")
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		httpx.Error(w, 400, "photo must be a multipart file up to 10 MiB")
		return
	}
	f, h, err := r.FormFile("photo")
	if err != nil {
		httpx.Error(w, 400, "photo field is required")
		return
	}
	defer f.Close()
	if !strings.HasPrefix(h.Header.Get("Content-Type"), "image/") {
		httpx.Error(w, 400, "only image uploads are allowed")
		return
	}
	key := fmt.Sprintf("doctors/%d/%d%s", id, time.Now().UnixNano(), filepath.Ext(h.Filename))
	url, err := i.s3.PutDoctorPhoto(r.Context(), key, h.Header.Get("Content-Type"), io.LimitReader(f, 10<<20))
	if err != nil {
		writeAction(w, nil, err, 500)
		return
	}
	v, err := i.crud.Update(r.Context(), "doctors", id, map[string]any{"photo": url})
	writeAction(w, v, err, 200)
}
func (i *initializer) getDoctorsOnMap(w http.ResponseWriter, r *http.Request) {
	var q struct {
		Lat      float64 `json:"lat"`
		Lon      float64 `json:"lon"`
		RadiusKm float64 `json:"radius_km"`
	}
	if err := httpx.Decode(r, &q); err != nil {
		httpx.Error(w, 400, "invalid request")
		return
	}
	v, err := i.onMap.Execute(r.Context(), q.Lat, q.Lon, q.RadiusKm)
	writeAction(w, v, err, 200)
}
func (i *initializer) filterReviews(w http.ResponseWriter, r *http.Request) {
	var q struct {
		CourseID *int64   `json:"course_id"`
		Lat      *float64 `json:"lat"`
		Lon      *float64 `json:"lon"`
		RadiusKm float64  `json:"radius_km"`
	}
	if err := httpx.Decode(r, &q); err != nil {
		httpx.Error(w, 400, "invalid request")
		return
	}
	if (q.Lat == nil) != (q.Lon == nil) {
		httpx.Error(w, 400, "lat and lon must be provided together")
		return
	}
	v, err := i.filter.Execute(r.Context(), q.CourseID, q.Lat, q.Lon, q.RadiusKm)
	writeAction(w, v, err, 200)
}
func decodeValues(r *http.Request) (map[string]any, error) {
	v := map[string]any{}
	return v, httpx.Decode(r, &v)
}
func hashPassword(v map[string]any) error {
	raw, ok := v["password"].(string)
	if !ok {
		return nil
	}
	if len(raw) < 10 {
		return fmt.Errorf("password must have at least 10 characters")
	}
	hash, err := password.Hash(raw)
	if err == nil {
		v["password"] = hash
	}
	return err
}
func idParam(r *http.Request) (int64, error) { return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64) }
func writeAction(w http.ResponseWriter, v any, err error, status int) {
	if err == nil {
		if v != nil {
			httpx.JSON(w, status, v)
		}
		return
	}
	if err == pgx.ErrNoRows {
		httpx.Error(w, 404, "not found")
		return
	}
	httpx.Error(w, 500, err.Error())
}
