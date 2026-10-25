package dal

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/dto"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Repository struct{ db postgres.DB }

func New(db postgres.DB) *Repository { return &Repository{db: db} }
func scan(row pgx.Row) (*dto.User, error) {
	var v dto.User
	err := row.Scan(&v.ID, &v.Email)
	if err != nil {
		return nil, postgres.Error(err)
	}
	return &v, nil
}
func (r *Repository) GetUser(ctx context.Context, id int64) (*dto.User, error) {
	return scan(r.db.QueryRow(ctx, `SELECT id, email FROM users WHERE id=$1`, id))
}
func (r *Repository) CreateUser(ctx context.Context, input dto.CreateInput) (*dto.User, error) {
	return scan(r.db.QueryRow(ctx, `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email`, input.Email, input.Password))
}
func (r *Repository) UpdateUser(ctx context.Context, id int64, input dto.UpdateInput) (*dto.User, error) {
	return scan(r.db.QueryRow(ctx, `UPDATE users SET email=COALESCE($2, email), password=COALESCE($3, password) WHERE id=$1 RETURNING id, email`, id, input.Email, input.Password))
}
func (r *Repository) ListUsers(ctx context.Context) ([]dto.User, error) {
	rows, err := r.db.Query(ctx, `SELECT id, email FROM users ORDER BY id`)
	if err != nil {
		return nil, postgres.Error(err)
	}
	return postgres.Collect(rows, scan)
}
func (r *Repository) DeleteUser(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	return postgres.DeleteResult(tag, err)
}
func (r *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, email).Scan(&exists)
	return exists, postgres.Error(err)
}
func (r *Repository) CreateInitialUser(ctx context.Context, email, hash string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO users(email,password) VALUES($1,$2) ON CONFLICT(email) DO NOTHING`, email, hash)
	return postgres.Error(err)
}
