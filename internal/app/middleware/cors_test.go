package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type originConfigStub struct {
	origins []string
}

func (s originConfigStub) GetAllowedOrigins() []string {
	return s.origins
}

func TestCORSMiddlewareAllowsConfiguredOriginWithCredentials(t *testing.T) {
	handler := CORSMiddleware(originConfigStub{origins: []string{"https://reviews.example.com"}})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/doctors/filter", nil)
	req.Header.Set("Origin", "https://reviews.example.com")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "https://reviews.example.com" {
		t.Errorf("Access-Control-Allow-Origin = %q", got)
	}
	if got := response.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("Access-Control-Allow-Credentials = %q", got)
	}
	if got := response.Header().Get("Access-Control-Allow-Headers"); got != "Content-Type, Authorization" {
		t.Errorf("Access-Control-Allow-Headers = %q", got)
	}
}

func TestCORSMiddlewareRequiresExactOriginForPreflight(t *testing.T) {
	handler := CORSMiddleware(originConfigStub{origins: []string{"https://reviews.example.com"}})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://reviews.example.com")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusForbidden)
	}
}

func TestCORSMiddlewareRejectsMalformedOriginForPreflight(t *testing.T) {
	handler := CORSMiddleware(originConfigStub{origins: []string{"https://reviews.example.com"}})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "\x7f")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusForbidden)
	}
}

func TestCORSMiddlewareAllowsServerToServerRequestWithoutOrigin(t *testing.T) {
	reachedHandler := false
	handler := CORSMiddleware(originConfigStub{origins: []string{"https://reviews.example.com"}})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		reachedHandler = true
		w.WriteHeader(http.StatusOK)
	}))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))

	if !reachedHandler {
		t.Fatal("handler was not called")
	}
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
}
