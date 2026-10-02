package validation

import (
	"errors"

	"github.com/nagorneva/nagorneva_reviews/internal/pkg/password"
)

var ErrPasswordTooShort = errors.New("password must have at least 10 characters")

func HashPassword(raw string) (string, error) {
	if len(raw) < 10 {
		return "", ErrPasswordTooShort
	}
	return password.Hash(raw)
}
