package dal

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/domain"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/geo"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Repository struct{ db postgres.DB }

func New(db postgres.DB) *Repository { return &Repository{db: db} }

func (d *Repository) Filter(ctx context.Context, courseID *int64, lat, lon *float64, radius float64) ([]domain.ReviewWithCourse, error) {
	south, north, west, east := -90.0, 90.0, -180.0, 180.0
	if lat != nil && lon != nil {
		south, north, west, east = geo.Bounds(*lat, *lon, radius)
	}
	rows, err := d.db.Query(ctx, `SELECT r.id,r.doctor_id,r.rating,r.comment,r.course_id,r.created_at,r.is_active,c.name,c.site_link FROM reviews r JOIN courses c ON c.id=r.course_id JOIN doctors d ON d.id=r.doctor_id JOIN cities city ON city.id=d.city_id WHERE r.is_active AND d.is_active AND city_id != 2 AND d.lat BETWEEN $5::float8 AND $6::float8 AND (($7::float8 <= $8::float8 AND d.lon BETWEEN $7::float8 AND $8::float8) OR ($7::float8 > $8::float8 AND (d.lon >= $7::float8 OR d.lon <= $8::float8))) AND ($1::bigint IS NULL OR r.course_id=$1) AND ($2::float8 IS NULL OR $4<=0 OR 6371*acos(least(1.0,cos(radians($2))*cos(radians(d.lat))*cos(radians(d.lon)-radians($3))+sin(radians($2))*sin(radians(d.lat))))<=$4) ORDER BY r.created_at DESC`, courseID, lat, lon, radius, south, north, west, east)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.ReviewWithCourse{}
	for rows.Next() {
		var r domain.ReviewWithCourse
		if err := rows.Scan(&r.ID, &r.DoctorID, &r.Rating, &r.Comment, &r.CourseID, &r.CreatedAt, &r.IsActive, &r.CourseName, &r.CourseSiteLink); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
