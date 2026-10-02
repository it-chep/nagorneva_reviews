// Package admin composes all admin actions into one application module.
package admin

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nagorneva/nagorneva_reviews/internal/config"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action"
	s3gateway "github.com/nagorneva/nagorneva_reviews/internal/module/admin/client/s3"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/auth"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/password"
)

// Module is the entry point for the admin domain.
type Module struct {
	Actions *action.Aggregator
}

// New initializes module-owned infrastructure and constructs admin actions.
func New(ctx context.Context, db *pgxpool.Pool, token auth.Service, cfg config.Config) (*Module, error) {
	if err := ensureInitialUser(ctx, db, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		return nil, err
	}
	s3Client, err := s3gateway.NewGateway(ctx, cfg)
	if err != nil {
		return nil, err
	}
	actions := action.NewAggregator(db, token, s3Client)
	return &Module{Actions: actions}, nil
}

func ensureInitialUser(ctx context.Context, db *pgxpool.Pool, email, rawPassword string) error {
	if email == "" || rawPassword == "" {
		return nil
	}
	var exists bool
	if err := db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, email).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil
	}
	hash, err := password.Hash(rawPassword)
	if err != nil {
		return err
	}
	_, err = db.Exec(ctx, `INSERT INTO users(email,password) VALUES($1,$2)`, email, hash)
	return err
}
