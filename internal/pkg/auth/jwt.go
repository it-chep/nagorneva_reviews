package auth

import (
	"errors"
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
