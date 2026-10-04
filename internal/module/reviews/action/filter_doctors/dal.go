package filter_doctors

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nagorneva/nagorneva_reviews/internal/domain"
)

type DAL struct{ db *pgxpool.Pool }

// Filter combines non-empty filter groups with AND. IDs inside one group are
// alternatives, so a doctor may match any city, specialty, or completed course
// supplied by the caller.
func (d DAL) Filter(ctx context.Context, cityIDs, specialtyIDs, courseIDs []int64, isActive, personalDataConsent *bool, reviewsSort ReviewsSort) ([]domain.FilteredDoctor, error) {
	rows, err := d.db.Query(ctx, `
WITH filtered_doctors AS (
SELECT
    d.id,
    d.name,
    CASE WHEN d.personal_data_consent THEN d.photo ELSE '' END AS photo,
    city.id AS city_id, city.name AS city_name, city.lat AS city_lat, city.lon AS city_lon,
    specialty.id AS specialty_id, specialty.name AS specialty_name,
    d.is_active, d.personal_data_consent,
    (SELECT count(*) FROM reviews review WHERE review.doctor_id = d.id AND review.is_active) AS reviews_count
FROM doctors d
JOIN cities city ON city.id = d.city_id
JOIN specialties specialty ON specialty.id = d.specialty_id
WHERE (COALESCE(cardinality($1::bigint[]), 0) = 0 OR d.city_id = ANY($1::bigint[]))
  AND (COALESCE(cardinality($2::bigint[]), 0) = 0 OR d.specialty_id = ANY($2::bigint[]))
  AND (
      COALESCE(cardinality($3::bigint[]), 0) = 0
      OR EXISTS (
          SELECT 1
          FROM doctor_courses completion
          WHERE completion.doctor_id = d.id AND completion.course_id = ANY($3::bigint[])
      )
  )
  AND ($4::boolean IS NULL OR d.is_active = $4)
  AND ($5::boolean IS NULL OR d.personal_data_consent = $5)
)
SELECT * FROM filtered_doctors
ORDER BY
  CASE WHEN $6 = 1 THEN reviews_count END DESC,
  CASE WHEN $6 = 2 THEN reviews_count END ASC,
  id ASC`, cityIDs, specialtyIDs, courseIDs, isActive, personalDataConsent, int32(reviewsSort))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	doctors := make([]domain.FilteredDoctor, 0)
	doctorIDs := make([]int64, 0)
	for rows.Next() {
		var doctor domain.FilteredDoctor
		if err := rows.Scan(
			&doctor.ID,
			&doctor.Name,
			&doctor.Photo,
			&doctor.City.ID,
			&doctor.City.Name,
			&doctor.City.Lat,
			&doctor.City.Lon,
			&doctor.Specialty.ID,
			&doctor.Specialty.Name,
			&doctor.IsActive,
			&doctor.PersonalDataConsent,
			&doctor.ReviewsCount,
		); err != nil {
			return nil, err
		}
		doctors = append(doctors, doctor)
		doctorIDs = append(doctorIDs, doctor.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	completedCourses, err := d.completedCourses(ctx, doctorIDs)
	if err != nil {
		return nil, err
	}
	for index := range doctors {
		doctors[index].CompletedCourses = completedCourses[doctors[index].ID]
		if doctors[index].CompletedCourses == nil {
			doctors[index].CompletedCourses = make([]domain.Course, 0)
		}
	}
	return doctors, nil
}

func (d DAL) completedCourses(ctx context.Context, doctorIDs []int64) (map[int64][]domain.Course, error) {
	result := make(map[int64][]domain.Course)
	if len(doctorIDs) == 0 {
		return result, nil
	}

	rows, err := d.db.Query(ctx, `
SELECT completion.doctor_id, course.id, course.name
FROM doctor_courses completion
JOIN courses course ON course.id = completion.course_id
WHERE completion.doctor_id = ANY($1::bigint[])
ORDER BY completion.doctor_id, course.id`, doctorIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var doctorID int64
		var course domain.Course
		if err := rows.Scan(&doctorID, &course.ID, &course.Name); err != nil {
			return nil, err
		}
		result[doctorID] = append(result[doctorID], course)
	}
	return result, rows.Err()
}
