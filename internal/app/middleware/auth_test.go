package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc"
)

type rejectingVerifier struct{}

func (rejectingVerifier) Verify(string) (int64, error) {
	return 0, errors.New("token must not be verified")
}

func TestAdminAuthDebugBypassesAuthorization(t *testing.T) {
	handler := AdminAuth(rejectingVerifier{}, true)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil))
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
}

func TestGRPCAdminAuthDebugBypassesAuthorization(t *testing.T) {
	called := false
	interceptor := GRPCAdminAuth(rejectingVerifier{}, true)
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/nagorneva_reviews.admin.v1.AdminService/CreateUser",
	}, func(context.Context, any) (any, error) {
		called = true
		return nil, nil
	})
	if err != nil {
		t.Fatalf("interceptor error = %v", err)
	}
	if !called {
		t.Fatal("handler was not called")
	}
}
