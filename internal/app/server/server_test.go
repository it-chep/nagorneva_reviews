package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
	reviewspb "github.com/nagorneva/nagorneva_reviews/gen/go/reviews"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc/status"
)

type tokens struct{}

func (tokens) Verify(token string) (int64, error) {
	if token != "valid" {
		return 0, errors.New("invalid token")
	}
	return 1, nil
}

type adminStub struct {
	adminpb.UnimplementedAdminServiceServer
}

func (adminStub) Login(context.Context, *adminpb.LoginRequest) (*adminpb.LoginResponse, error) {
	return &adminpb.LoginResponse{AccessToken: "valid", TokenType: "Bearer"}, nil
}
func (adminStub) ListUsers(ctx context.Context, _ *adminpb.ListUsersRequest) (*adminpb.ListUsersResponse, error) {
	return &adminpb.ListUsersResponse{Users: []*adminpb.User{{Id: 1, Email: "admin@example.test"}}}, nil
}
func (adminStub) GetUser(_ context.Context, req *adminpb.GetUserRequest) (*adminpb.GetUserResponse, error) {
	if req.Id == 500 {
		panic("test panic")
	}
	return nil, status.Error(codes.NotFound, "not found")
}
func (adminStub) UpdateDoctor(_ context.Context, req *adminpb.UpdateDoctorRequest) (*adminpb.UpdateDoctorResponse, error) {
	if req.IsActive == nil || *req.IsActive || req.Name == nil || *req.Name != "" || req.FullName != nil {
		return nil, status.Error(codes.InvalidArgument, "optional field presence lost")
	}
	return &adminpb.UpdateDoctorResponse{Doctor: &adminpb.Doctor{Id: req.Id}}, nil
}
func (adminStub) UploadDoctorPhoto(_ context.Context, req *adminpb.UploadDoctorPhotoRequest) (*adminpb.UploadDoctorPhotoResponse, error) {
	return &adminpb.UploadDoctorPhotoResponse{Doctor: &adminpb.Doctor{Id: req.Id, ReviewsCount: int64(len(req.Photo))}}, nil
}
func (adminStub) UploadDoctorPhotoMultipart(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

type reviewsStub struct {
	reviewspb.UnimplementedReviewsServiceServer
}

func (reviewsStub) FilterDoctors(ctx context.Context, _ *reviewspb.FilterDoctorsRequest) (*reviewspb.FilterDoctorsResponse, error) {
	return &reviewspb.FilterDoctorsResponse{DoctorCount: 2}, nil
}

func runningServer(t *testing.T, debug bool) (string, *grpc.ClientConn) {
	t.Helper()
	rpcListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	httpListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		_ = rpcListener.Close()
		t.Fatal(err)
	}
	rpcAddr, httpAddr := rpcListener.Addr().String(), httpListener.Addr().String()
	_ = rpcListener.Close()
	_ = httpListener.Close()
	cfg := Config{HTTPAddr: httpAddr, GRPCAddr: rpcAddr, Token: tokens{}, Debug: debug, AllowedOrigins: []string{"https://admin.example.test"}}
	s, err := New(context.Background(), cfg, adminStub{}, reviewsStub{})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()
	t.Cleanup(func() {
		cancel()
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("server shutdown: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Error("server did not stop")
		}
		_ = rpcListener.Close()
		_ = httpListener.Close()
	})
	conn, err := grpc.NewClient("passthrough:///"+rpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithNoProxy())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	select {
	case <-s.services.ready:
	case err := <-done:
		t.Fatalf("server startup: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("server startup timed out")
	}
	return "http://" + httpAddr, conn
}

func TestGatewayAndNativeGRPC(t *testing.T) {
	base, conn := runningServer(t, false)
	client := adminpb.NewAdminServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := client.ListUsers(ctx, &adminpb.ListUsersRequest{}); status.Code(err) != codes.Unauthenticated {
		t.Fatalf("native auth: %v", err)
	}
	authed := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer valid")
	users, err := client.ListUsers(authed, &adminpb.ListUsersRequest{})
	if err != nil || len(users.GetUsers()) != 1 {
		t.Fatalf("native list: %v %v", users, err)
	}
	cases := []struct {
		name, method, path, body, token string
		code                            int
	}{
		{"missing token", "GET", "/api/v1/admin/users", "", "", 401},
		{"invalid token", "GET", "/api/v1/admin/users", "", "bad", 401},
		{"gateway auth", "GET", "/api/v1/admin/users", "", "valid", 200},
		{"trailing slash", "GET", "/api/v1/admin/users/", "", "valid", 200},
		{"login exemption", "POST", "/api/v1/admin/login", `{"email":"a","password":"b"}`, "", 200},
		{"public", "POST", "/api/v1/doctors/filter", `{}`, "", 200},
		{"partial update", "PATCH", "/api/v1/admin/doctors/7", `{"is_active":false,"name":""}`, "valid", 200},
		{"not found", "GET", "/api/v1/admin/users/404", "", "valid", 404},
		{"panic recovery", "GET", "/api/v1/admin/users/500", "", "valid", 500},
		{"health", "GET", "/healthz", "", "", 200},
		{"swagger", "GET", "/swagger.json", "", "", 200},
	}
	httpClient := &http.Client{Timeout: 5 * time.Second}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, base+tc.path, strings.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			response, err := httpClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = response.Body.Close() }()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}
			if response.StatusCode != tc.code {
				t.Fatalf("status=%d, body=%s", response.StatusCode, body)
			}
			if tc.name == "gateway auth" && !bytes.Contains(body, []byte("admin@example.test")) {
				t.Fatalf("missing result: %s", body)
			}
			if tc.name == "login exemption" && !bytes.Contains(body, []byte(`"accessToken"`)) {
				t.Fatalf("JSON field naming changed: %s", body)
			}
			if tc.name == "swagger" && (!bytes.Contains(body, []byte("/api/v1/admin/users")) || !bytes.Contains(body, []byte("/api/v1/doctors/filter"))) {
				t.Fatal("incomplete combined OpenAPI")
			}
		})
	}
}

func TestGatewayPhotoAndCORS(t *testing.T) {
	base, _ := runningServer(t, false)
	httpClient := &http.Client{Timeout: 10 * time.Second}
	for _, tc := range []struct {
		kind, token string
		code        int
	}{
		{"multipart/form-data; boundary=test", "", 401},
		{"multipart/form-data; boundary=test", "valid", 202},
		{"application/json", "", 401},
		{"application/json", "valid", 200},
	} {
		body := []byte(`{"photo":"YQ==","filename":"photo.png","content_type":"image/png"}`)
		req, err := http.NewRequest("POST", base+"/api/v1/admin/doctors/7/photo", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", tc.kind)
		if tc.token != "" {
			req.Header.Set("Authorization", "Bearer "+tc.token)
		}
		response, err := httpClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		_ = response.Body.Close()
		if response.StatusCode != tc.code {
			t.Errorf("%s token=%s: status=%d", tc.kind, tc.token, response.StatusCode)
		}
	}
	// Exercise payloads exceeding gRPC's default 4 MiB receive limit.
	payload, err := json.Marshal(map[string]any{"photo": make([]byte, 5<<20), "filename": "large.png", "content_type": "image/png"})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", base+"/api/v1/admin/doctors/7/photo", bytes.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid")
	response, err := httpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 || !bytes.Contains(body, []byte(`"5242880"`)) {
		t.Fatalf("large photo: %d %s", response.StatusCode, body)
	}
	req, err = http.NewRequest("OPTIONS", base+"/api/v1/admin/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Origin", "https://admin.example.test")
	response, err = httpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = response.Body.Close()
	if response.StatusCode != 204 || response.Header.Get("Access-Control-Allow-Origin") != "https://admin.example.test" {
		t.Fatalf("CORS: %d %v", response.StatusCode, response.Header)
	}
}

func TestDebugBypassesBothTransports(t *testing.T) {
	base, conn := runningServer(t, true)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := adminpb.NewAdminServiceClient(conn).ListUsers(ctx, &adminpb.ListUsersRequest{}); err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Get(base + "/api/v1/admin/users")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != 200 {
		t.Fatal(response.StatusCode)
	}
}

func TestRunDoesNotStartWithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s, err := New(context.Background(), Config{HTTPAddr: ":8080", GRPCAddr: ":7002", Token: tokens{}}, adminStub{}, reviewsStub{})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Run(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
}

func TestMetadataHeaderCannotReplaceAuthorization(t *testing.T) {
	base, _ := runningServer(t, false)
	req, err := http.NewRequest("GET", base+"/api/v1/admin/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Grpc-Metadata-Authorization", "Bearer valid")
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unexpected auth bypass: %d", response.StatusCode)
	}
}
