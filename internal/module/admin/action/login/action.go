package login

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/password"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Action struct{ dal DAL }

func New(db *pgxpool.Pool) Action { return Action{DAL{db: db}} }
func (a Action) Execute(ctx context.Context, email, rawPassword string) (int64, error) {
	id, hash, err := a.dal.FindByEmail(ctx, email)
	if err != nil {
		return 0, ErrInvalidCredentials
	}
	if password.Check(hash, rawPassword) != nil {
		return 0, ErrInvalidCredentials
	}
	return id, nil
}
