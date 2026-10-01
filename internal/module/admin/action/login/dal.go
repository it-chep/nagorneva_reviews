package login

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DAL struct{ db *pgxpool.Pool }

func (d DAL) FindByEmail(ctx context.Context, email string) (int64, string, error) {
	var id int64
	var hash string
	err := d.db.QueryRow(ctx, `SELECT id,password FROM users WHERE email=$1`, email).Scan(&id, &hash)
	return id, hash, err
}
