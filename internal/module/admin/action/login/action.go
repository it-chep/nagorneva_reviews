package login

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/auth"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/password"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Action struct {
	dal   DAL
	token auth.Service
}

func New(db *pgxpool.Pool, token auth.Service) *Action {
	return &Action{dal: DAL{db: db}, token: token}
}

func (a *Action) Execute(ctx context.Context, email, rawPassword string) (string, error) {
	id, hash, err := a.dal.FindByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if password.Check(hash, rawPassword) != nil {
		return "", ErrInvalidCredentials
	}
	return a.token.Issue(id)
}
