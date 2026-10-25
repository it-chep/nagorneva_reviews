package login

import (
	"context"
	"errors"

	"github.com/nagorneva/nagorneva_reviews/internal/pkg/domain_error"

	"github.com/nagorneva/nagorneva_reviews/internal/pkg/password"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Storage interface {
	FindByEmail(context.Context, string) (int64, string, error)
}
type TokenIssuer interface{ Issue(int64) (string, error) }
type Action struct {
	dal   Storage
	token TokenIssuer
}

func New(storage Storage, token TokenIssuer) *Action { return &Action{dal: storage, token: token} }

func (a *Action) Execute(ctx context.Context, email, rawPassword string) (string, error) {
	id, hash, err := a.dal.FindByEmail(ctx, email)
	if errors.Is(err, domain_error.ErrNotFound) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}
	if password.Check(hash, rawPassword) != nil {
		return "", ErrInvalidCredentials
	}
	return a.token.Issue(id)
}
