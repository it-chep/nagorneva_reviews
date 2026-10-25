package dal

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/dto"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Repository struct{ db postgres.DB }

func New(db postgres.DB) *Repository { return &Repository{db: db} }

const selection = `d.id, d.name, d.full_name, d.city_id, d.specialty_id, d.lat, d.lon, d.personal_data_consent, d.is_active, d.photo, c.name, c.lat, c.lon, s.name, (SELECT count(*) FROM reviews review WHERE review.doctor_id=d.id)`
const joins = `JOIN cities c ON c.id=d.city_id JOIN specialties s ON s.id=d.specialty_id`

func scan(row pgx.Row) (*dto.Doctor, error) {
	var v dto.Doctor
	err := row.Scan(&v.ID, &v.Name, &v.FullName, &v.CityID, &v.SpecialtyID, &v.Lat, &v.Lon, &v.PersonalDataConsent, &v.IsActive, &v.Photo, &v.CityName, &v.CityLat, &v.CityLon, &v.SpecialtyName, &v.ReviewsCount)
	if err != nil {
		return nil, postgres.Error(err)
	}
	return &v, nil
}
func (r *Repository) GetDoctor(ctx context.Context, id int64) (*dto.Doctor, error) {
	return scan(r.db.QueryRow(ctx, `SELECT `+selection+` FROM doctors d `+joins+` WHERE d.id=$1`, id))
}
func (r *Repository) ListDoctors(ctx context.Context) ([]dto.Doctor, error) {
	rows, err := r.db.Query(ctx, `SELECT `+selection+` FROM doctors d `+joins+` ORDER BY d.id`)
	if err != nil {
		return nil, postgres.Error(err)
	}
	return postgres.Collect(rows, scan)
}
func (r *Repository) CreateDoctor(ctx context.Context, input dto.CreateInput) (*dto.Doctor, error) {
	return scan(r.db.QueryRow(ctx, `WITH changed AS (INSERT INTO doctors (name, full_name, city_id, specialty_id, lat, lon, personal_data_consent, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *) SELECT `+selection+` FROM changed d `+joins, input.Name, input.FullName, input.CityID, input.SpecialtyID, input.Lat, input.Lon, input.PersonalDataConsent, input.IsActive))
}
func (r *Repository) UpdateDoctor(ctx context.Context, id int64, input dto.UpdateInput) (*dto.Doctor, error) {
	return scan(r.db.QueryRow(ctx, `WITH changed AS (UPDATE doctors SET name=COALESCE($2, name), full_name=COALESCE($3, full_name), city_id=COALESCE($4, city_id), specialty_id=COALESCE($5, specialty_id), lat=COALESCE($6, lat), lon=COALESCE($7, lon), personal_data_consent=COALESCE($8, personal_data_consent), is_active=COALESCE($9, is_active), updated_at=now() WHERE id=$1 RETURNING *) SELECT `+selection+` FROM changed d `+joins, id, input.Name, input.FullName, input.CityID, input.SpecialtyID, input.Lat, input.Lon, input.PersonalDataConsent, input.IsActive))
}
func (r *Repository) SetPhoto(ctx context.Context, id int64, photo string) (*dto.Doctor, error) {
	return scan(r.db.QueryRow(ctx, `WITH changed AS (UPDATE doctors SET photo=$2, updated_at=now() WHERE id=$1 RETURNING *) SELECT `+selection+` FROM changed d `+joins, id, photo))
}
func (r *Repository) DeleteDoctor(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM doctors WHERE id=$1`, id)
	return postgres.DeleteResult(tag, err)
}
