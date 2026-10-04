// Package server owns the application's HTTP grpc-gateway and native gRPC
// transports. It uses the same Clay setup as medblogers_base so gateway calls
// pass through the generated service descriptors and gRPC interceptors.
package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
	reviewspb "github.com/nagorneva/nagorneva_reviews/gen/go/reviews"
	healthv1 "github.com/nagorneva/nagorneva_reviews/internal/app/api/health/v1"
	appmiddleware "github.com/nagorneva/nagorneva_reviews/internal/app/middleware"
	clayserver "github.com/not-for-prod/clay/server"
	"github.com/not-for-prod/clay/transport"
)

const shutdownTimeout = 15 * time.Second

// Config configures the public HTTP gateway and native gRPC listener.
type Config struct {
	HTTPAddr string
	GRPCAddr string
	Token    appmiddleware.TokenVerifier
	Debug    bool
}

// AdminService is the generated gRPC service plus the browser multipart
// endpoint that grpc-gateway cannot represent.
type AdminService interface {
	adminpb.AdminServiceServer
	UploadDoctorPhotoMultipart(http.ResponseWriter, *http.Request)
}

// Server runs Clay's generated-descriptor transport. Clay registers each
// service both in native gRPC and grpc-gateway, applying unary middleware to
// both paths just as medblogers_base does.
type Server struct {
	server      *clayserver.Server
	controllers []transport.ServiceDesc
	httpAddr    string
	grpcAddr    string
}

// New composes the Chi routes that are not generated from protobuf (health
// and multipart upload) and leaves protobuf HTTP/gRPC registration to Clay.
func New(_ context.Context, cfg Config, adminService AdminService, reviewsService reviewspb.ReviewsServiceServer) (*Server, error) {
	httpPort, err := addressPort("HTTP_ADDR", cfg.HTTPAddr)
	if err != nil {
		return nil, err
	}
	grpcPort, err := addressPort("GRPC_ADDR", cfg.GRPCAddr)
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()
	// Clay applies WithHTTPMiddlewares during Run. This project has explicit
	// non-protobuf routes, so the middleware must be attached before those
	// routes are declared (Chi rejects middleware added afterwards).
	router.Use(
		// grpc-gateway treats `/resource/` as `/resource/{id}` with an empty
		// value. Normalize it before gateway routing so list endpoints work
		// consistently with and without a trailing slash.
		trimTrailingSlash,
		chimiddleware.Recoverer,
		chimiddleware.RequestID,
		chimiddleware.RealIP,
		chimiddleware.Timeout(30*time.Second),
		appmiddleware.CORS,
	)
	router.Get("/healthz", healthv1.Get)
	router.With(appmiddleware.AdminAuth(cfg.Token, cfg.Debug)).Post("/api/v1/admin/doctors/{id}/photo", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/") {
			adminService.UploadDoctorPhotoMultipart(w, r)
			return
		}
		http.Error(w, `{"error":"multipart/form-data is required"}`, http.StatusUnsupportedMediaType)
	})

	clay := clayserver.NewServer(
		grpcPort,
		clayserver.WithHTTPPort(httpPort),
		clayserver.WithHTTPMux(router),
		clayserver.WithGRPCUnaryMiddlewares(appmiddleware.GRPCAdminAuth(cfg.Token, cfg.Debug)),
		clayserver.WithRuntimeServeMuxOpts(runtime.WithIncomingHeaderMatcher(incomingHeaderMatcher)),
	)

	return &Server{
		server: clay,
		controllers: []transport.ServiceDesc{
			adminpb.NewAdminServiceServiceDesc(adminService),
			reviewspb.NewReviewsServiceServiceDesc(reviewsService),
		},
		httpAddr: cfg.HTTPAddr,
		grpcAddr: cfg.GRPCAddr,
	}, nil
}

func incomingHeaderMatcher(key string) (string, bool) {
	if strings.EqualFold(key, "Authorization") {
		return "authorization", true
	}
	return runtime.DefaultHeaderMatcher(key)
}

// trimTrailingSlash updates both URL.Path and Chi's route context. Chi's own
// StripSlashes only changes the route context, while Clay's root-mounted
// gateway computes its child route from URL.Path as well.
func trimTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 1 && strings.HasSuffix(r.URL.Path, "/") {
			path := strings.TrimRight(r.URL.Path, "/")
			r.URL.Path = path
			r.URL.RawPath = ""
			if routeContext := chi.RouteContext(r.Context()); routeContext != nil {
				routeContext.RoutePath = path
			}
		}
		next.ServeHTTP(w, r)
	})
}

func addressPort(name, address string) (int, error) {
	_, rawPort, err := net.SplitHostPort(address)
	if err != nil {
		return 0, fmt.Errorf("%s %q must be host:port: %w", name, address, err)
	}
	port, err := strconv.Atoi(rawPort)
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("%s %q must contain a port from 1 to 65535", name, address)
	}
	return port, nil
}

// HTTPAddr returns the configured grpc-gateway listener address.
func (s *Server) HTTPAddr() string { return s.httpAddr }

// GRPCAddr returns the configured native gRPC listener address.
func (s *Server) GRPCAddr() string { return s.grpcAddr }

// Run starts both transports. On cancellation it follows Clay's graceful
// shutdown path, matching medblogers_base's application lifecycle.
func (s *Server) Run(ctx context.Context) error {
	if s == nil || s.server == nil {
		return fmt.Errorf("application servers are not initialized")
	}

	errCh := make(chan error, 1)
	go func() { errCh <- s.server.Run(s.controllers...) }()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := s.Stop(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown: %w", err)
		}
		return nil
	}
}

// Stop gracefully terminates the HTTP and gRPC transports.
func (s *Server) Stop(ctx context.Context) error {
	if s == nil || s.server == nil {
		return nil
	}
	return s.server.Stop(ctx)
}
