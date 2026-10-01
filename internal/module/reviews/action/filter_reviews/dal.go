package filter_reviews

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nagorneva/nagorneva_reviews/internal/domain"
)

type DAL struct{ db *pgxpool.Pool }

func (d DAL) Filter(ctx context.Context, courseID *int64, lat, lon *float64, radius float64) ([]domain.ReviewWithCourse, error) {
	rows, err := d.db.Query(ctx, `SELECT r.id,r.doctor_id,r.rating,r.comment,r.course_id,r.created_at,r.is_active,c.name FROM reviews r JOIN courses c ON c.id=r.course_id JOIN doctors d ON d.id=r.doctor_id WHERE r.is_active AND d.is_active AND ($1::bigint IS NULL OR r.course_id=$1) AND ($2::float8 IS NULL OR $4<=0 OR 6371*acos(least(1.0,cos(radians($2))*cos(radians(d.lat))*cos(radians(d.lon)-radians($3))+sin(radians($2))*sin(radians(d.lat))))<=$4) ORDER BY r.created_at DESC`, courseID, lat, lon, radius)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.ReviewWithCourse{}
	for rows.Next() {
		var r domain.ReviewWithCourse
		if err := rows.Scan(&r.ID, &r.DoctorID, &r.Rating, &r.Comment, &r.CourseID, &r.CreatedAt, &r.IsActive, &r.CourseName); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
