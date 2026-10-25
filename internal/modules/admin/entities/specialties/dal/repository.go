package dal

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties/dto"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Repository struct{ db postgres.DB }

func New(db postgres.DB) *Repository { return &Repository{db: db} }
func scan(row pgx.Row) (*dto.Specialty, error) {
	var v dto.Specialty
	err := row.Scan(&v.ID, &v.Name)
	if err != nil {
		return nil, postgres.Error(err)
	}
	return &v, nil
}
func (r *Repository) GetSpecialty(ctx context.Context, id int64) (*dto.Specialty, error) {
	return scan(r.db.QueryRow(ctx, `SELECT id, name FROM specialties WHERE id=$1`, id))
}
func (r *Repository) CreateSpecialty(ctx context.Context, input dto.CreateInput) (*dto.Specialty, error) {
	return scan(r.db.QueryRow(ctx, `INSERT INTO specialties (name) VALUES ($1) RETURNING id, name`, input.Name))
}
func (r *Repository) UpdateSpecialty(ctx context.Context, id int64, input dto.UpdateInput) (*dto.Specialty, error) {
	return scan(r.db.QueryRow(ctx, `UPDATE specialties SET name=COALESCE($2, name) WHERE id=$1 RETURNING id, name`, id, input.Name))
}
func (r *Repository) ListSpecialties(ctx context.Context) ([]dto.Specialty, error) {
	rows, err := r.db.Query(ctx, `SELECT entity.id, entity.name, (SELECT count(*) FROM doctors related WHERE related.specialty_id = entity.id) FROM specialties entity ORDER BY entity.id`)
	if err != nil {
		return nil, postgres.Error(err)
	}
	return postgres.Collect(rows, func(row pgx.Row) (*dto.Specialty, error) {
		var v dto.Specialty
		err := row.Scan(&v.ID, &v.Name, &v.DoctorsCount)
		return &v, err
	})
}
func (r *Repository) DeleteSpecialty(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM specialties WHERE id=$1`, id)
	return postgres.DeleteResult(tag, err)
}
