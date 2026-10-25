package password

import (
	"errors"

	passwordhash "github.com/nagorneva/nagorneva_reviews/internal/pkg/password"
)

var ErrPasswordTooShort = errors.New("password must have at least 10 characters")

func Hash(raw string) (string, error) {
	if len(raw) < 10 {
		return "", ErrPasswordTooShort
	}
	return passwordhash.Hash(raw)
}
