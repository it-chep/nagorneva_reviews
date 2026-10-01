package get_doctors_on_map

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nagorneva/nagorneva_reviews/internal/domain"
)

type DAL struct{ db *pgxpool.Pool }

func (d DAL) Get(ctx context.Context, lat, lon, radius float64) ([]domain.DoctorOnMap, error) {
	// Equirectangular distance is adequate for UI filtering; DB keeps independent map coordinates.
	rows, err := d.db.Query(ctx, `SELECT d.id,d.name,CASE WHEN d.personal_data_consent THEN d.full_name ELSE '' END,CASE WHEN d.personal_data_consent THEN d.photo ELSE '' END,d.personal_data_consent,d.city_id,d.specialty_id,d.lat,d.lon,d.is_active,c.name FROM doctors d JOIN cities c ON c.id=d.city_id WHERE d.is_active AND ($3 <= 0 OR 6371 * acos(least(1.0, cos(radians($1))*cos(radians(d.lat))*cos(radians(d.lon)-radians($2))+sin(radians($1))*sin(radians(d.lat)))) <= $3) ORDER BY d.id`, lat, lon, radius)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.DoctorOnMap{}
	for rows.Next() {
		var v domain.DoctorOnMap
		if err := rows.Scan(&v.ID, &v.Name, &v.FullName, &v.Photo, &v.PersonalDataConsent, &v.CityID, &v.SpecialtyID, &v.Lat, &v.Lon, &v.IsActive, &v.City); err != nil {
			return nil, err
		}
		reviews, err := d.reviews(ctx, v.ID)
		if err != nil {
			return nil, err
		}
		v.Reviews = reviews
		out = append(out, v)
	}
	return out, rows.Err()
}
func (d DAL) reviews(ctx context.Context, doctorID int64) ([]domain.ReviewWithCourse, error) {
	rows, err := d.db.Query(ctx, `SELECT r.id,r.doctor_id,r.rating,r.comment,r.course_id,r.created_at,r.is_active,c.name FROM reviews r JOIN courses c ON c.id=r.course_id WHERE r.doctor_id=$1 AND r.is_active ORDER BY r.created_at DESC`, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []domain.ReviewWithCourse{}
	for rows.Next() {
		var r domain.ReviewWithCourse
		if err := rows.Scan(&r.ID, &r.DoctorID, &r.Rating, &r.Comment, &r.CourseID, &r.CreatedAt, &r.IsActive, &r.CourseName); err != nil {
			return nil, fmt.Errorf("scan review: %w", err)
		}
		result = append(result, r)
	}
	return result, rows.Err()
}
