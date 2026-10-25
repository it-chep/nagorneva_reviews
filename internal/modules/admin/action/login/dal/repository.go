package dal

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Repository struct{ db postgres.DB }

func New(db postgres.DB) *Repository { return &Repository{db: db} }

func (d *Repository) FindByEmail(ctx context.Context, email string) (int64, string, error) {
	var id int64
	var hash string
	err := d.db.QueryRow(ctx, `SELECT id,password FROM users WHERE email=$1`, email).Scan(&id, &hash)
	return id, hash, postgres.Error(err)
}
