// Package server runs Clay with generated gRPC and grpc-gateway descriptors.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
)

const shutdownTimeout = 15 * time.Second

// Photo requests contain up to 10 MiB of binary data plus protobuf metadata.
const maxMessageSize = 12 << 20

type Config struct {
	HTTPAddr, GRPCAddr string
	AllowedOrigins     []string
	Token              appmiddleware.TokenVerifier
	Debug              bool
}

func (c Config) GetAllowedOrigins() []string { return c.AllowedOrigins }

type AdminService interface {
	adminpb.AdminServiceServer
	UploadDoctorPhotoMultipart(http.ResponseWriter, *http.Request)
}

// registeredServices keeps gateway access for the multipart/JSON shared URL,
// and signals when Clay has finished initializing both transports.
type registeredServices struct {
	*transport.CompoundServiceDesc
	gateway *runtime.ServeMux
	rpc     *grpc.Server
	ready   chan struct{}
}

func (d *registeredServices) RegisterHTTP(ctx context.Context, mux *runtime.ServeMux) error {
	d.gateway = mux
	return d.CompoundServiceDesc.RegisterHTTP(ctx, mux)
}
func (d *registeredServices) RegisterGRPC(rpc *grpc.Server) {
	d.CompoundServiceDesc.RegisterGRPC(rpc)
	d.rpc = rpc
	close(d.ready)
}

type Server struct {
	cfg      Config
	clay     *clayserver.Server
	services *registeredServices
	admin    AdminService
}

func New(_ context.Context, cfg Config, admin AdminService, reviews reviewspb.ReviewsServiceServer) (*Server, error) {
	httpPort, err := addressPort("HTTP_ADDR", cfg.HTTPAddr)
	if err != nil {
		return nil, err
	}
	grpcPort, err := addressPort("GRPC_ADDR", cfg.GRPCAddr)
	if err != nil {
		return nil, err
	}
	if !cfg.Debug && cfg.Token == nil {
		return nil, fmt.Errorf("admin token verifier is required")
	}
	s := &Server{cfg: cfg, admin: admin, services: &registeredServices{
		CompoundServiceDesc: transport.NewCompoundServiceDesc(
			adminpb.NewAdminServiceServiceDesc(admin),
			reviewspb.NewReviewsServiceServiceDesc(reviews),
		),
		ready: make(chan struct{}),
	}}
	s.clay = clayserver.NewServer(grpcPort,
		clayserver.WithHTTPPort(httpPort),
		clayserver.WithHTTPMux(s.router()),
		clayserver.WithGRPCUnaryMiddlewares(recoverRPC, appmiddleware.GRPCAdminAuth(cfg.Token, cfg.Debug)),
		clayserver.WithGRPCOpts(grpc.MaxRecvMsgSize(maxMessageSize), grpc.MaxSendMsgSize(maxMessageSize)),
		clayserver.WithRuntimeServeMuxOpts(runtime.WithIncomingHeaderMatcher(incomingHeaderMatcher)),
	)
	return s, nil
}

func recoverRPC(ctx context.Context, req any, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (response any, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			slog.ErrorContext(ctx, "gRPC handler panic", "method", info.FullMethod, "panic", recovered)
			response = nil
			err = status.Error(codes.Internal, "internal server error")
		}
	}()
	return next(ctx, req)
}

func (s *Server) HTTPAddr() string { return s.cfg.HTTPAddr }
func (s *Server) GRPCAddr() string { return s.cfg.GRPCAddr }

// Run delegates listeners, gateway registration and gRPC registration to Clay.
func (s *Server) Run(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	results := make(chan error, 1)
	go func() { results <- s.clay.Run(s.services) }()
	// Clay initializes its server fields in Run. Wait for registration before
	// calling Stop so cancellation cannot race those writes.
	select {
	case err := <-results:
		return err
	case <-s.services.ready:
	}
	var runErr error
	completed := false
	select {
	case runErr = <-results:
		completed = true
	case <-ctx.Done():
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	stopped := make(chan error, 1)
	go func() { stopped <- s.clay.Stop(shutdownCtx) }()
	var stopErr error
	select {
	case stopErr = <-stopped:
	case <-shutdownCtx.Done():
		s.services.rpc.Stop()
		stopErr = shutdownCtx.Err()
	}
	if !completed {
		select {
		case runErr = <-results:
		case <-shutdownCtx.Done():
			runErr = shutdownCtx.Err()
		}
	}
	if errors.Is(runErr, http.ErrServerClosed) || errors.Is(runErr, grpc.ErrServerStopped) {
		runErr = nil
	}
	return errors.Join(runErr, stopErr)
}

func (s *Server) router() *chi.Mux {
	router := chi.NewRouter()
	router.Use(trimTrailingSlash, chimiddleware.Recoverer, chimiddleware.RequestID,
		chimiddleware.RealIP, chimiddleware.Timeout(30*time.Second), appmiddleware.CORSMiddleware(s.cfg))
	router.Get("/healthz", healthv1.Get)
	multipart := appmiddleware.AdminAuth(s.cfg.Token, s.cfg.Debug)(http.HandlerFunc(s.admin.UploadDoctorPhotoMultipart))
	router.Post("/api/v1/admin/doctors/{id}/photo", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/") {
			multipart.ServeHTTP(w, r)
			return
		}
		// JSON/base64 requests for this same RPC must reach the generated gateway.
		s.services.gateway.ServeHTTP(w, r)
	})
	return router
}

func incomingHeaderMatcher(key string) (string, bool) {
	// runtime forwards Authorization itself. Do not append a second copy or
	// accept an alternative credential supplied through a metadata header.
	if strings.EqualFold(key, "Grpc-Metadata-Authorization") {
		return "", false
	}
	return runtime.DefaultHeaderMatcher(key)
}

func trimTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 1 && strings.HasSuffix(r.URL.Path, "/") && !strings.HasPrefix(r.URL.Path, "/docs/") {
			r.URL.Path = strings.TrimRight(r.URL.Path, "/")
			r.URL.RawPath = ""
			if routeContext := chi.RouteContext(r.Context()); routeContext != nil {
				routeContext.RoutePath = r.URL.Path
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
		return 0, fmt.Errorf("%s %q has an invalid port", name, address)
	}
	return port, nil
}
