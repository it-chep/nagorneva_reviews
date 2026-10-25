package dal

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities/dto"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Repository struct{ db postgres.DB }

func New(db postgres.DB) *Repository { return &Repository{db: db} }
func scan(row pgx.Row) (*dto.City, error) {
	var v dto.City
	err := row.Scan(&v.ID, &v.Name, &v.Lat, &v.Lon)
	if err != nil {
		return nil, postgres.Error(err)
	}
	return &v, nil
}
func (r *Repository) GetCity(ctx context.Context, id int64) (*dto.City, error) {
	return scan(r.db.QueryRow(ctx, `SELECT id, name, lat, lon FROM cities WHERE id=$1`, id))
}
func (r *Repository) CreateCity(ctx context.Context, input dto.CreateInput) (*dto.City, error) {
	return scan(r.db.QueryRow(ctx, `INSERT INTO cities (name, lat, lon) VALUES ($1, $2, $3) RETURNING id, name, lat, lon`, input.Name, input.Lat, input.Lon))
}
func (r *Repository) UpdateCity(ctx context.Context, id int64, input dto.UpdateInput) (*dto.City, error) {
	return scan(r.db.QueryRow(ctx, `UPDATE cities SET name=COALESCE($2, name), lat=COALESCE($3, lat), lon=COALESCE($4, lon) WHERE id=$1 RETURNING id, name, lat, lon`, id, input.Name, input.Lat, input.Lon))
}
func (r *Repository) ListCities(ctx context.Context) ([]dto.City, error) {
	rows, err := r.db.Query(ctx, `SELECT entity.id, entity.name, entity.lat, entity.lon, (SELECT count(*) FROM doctors related WHERE related.city_id = entity.id) FROM cities entity ORDER BY entity.id`)
	if err != nil {
		return nil, postgres.Error(err)
	}
	return postgres.Collect(rows, func(row pgx.Row) (*dto.City, error) {
		var v dto.City
		err := row.Scan(&v.ID, &v.Name, &v.Lat, &v.Lon, &v.DoctorsCount)
		return &v, err
	})
}
func (r *Repository) DeleteCity(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM cities WHERE id=$1`, id)
	return postgres.DeleteResult(tag, err)
}
