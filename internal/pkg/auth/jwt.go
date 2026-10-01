package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	UserID int64 `json:"uid"`
	jwt.RegisteredClaims
}
type Service struct {
	secret []byte
	ttl    time.Duration
}

func New(secret string, ttl time.Duration) Service { return Service{[]byte(secret), ttl} }
func (s Service) Issue(userID int64) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims{userID, jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)), IssuedAt: jwt.NewNumericDate(time.Now())}}).SignedString(s.secret)
}
func (s Service) Verify(raw string) (int64, error) {
	var c claims
	token, err := jwt.ParseWithClaims(raw, &c, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid || c.UserID == 0 {
		return 0, errors.New("invalid token")
	}
	return c.UserID, nil
}

type contextKey struct{}

func UserID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(contextKey{}).(int64)
	return id, ok
}
func (s Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			http.Error(w, `{"error":"authorization required"}`, http.StatusUnauthorized)
			return
		}
		id, err := s.Verify(strings.TrimPrefix(h, "Bearer "))
		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKey{}, id)))
	})
}
