package dal

import (
	"context"

	"github.com/jackc/pgx/v5"
	coursedto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/dto"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/domain_error"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

func (r *Repository) GetDoctorCourses(ctx context.Context, doctorID int64) ([]coursedto.Course, error) {
	var id int64
	if err := r.db.QueryRow(ctx, `SELECT id FROM doctors WHERE id=$1`, doctorID).Scan(&id); err != nil {
		return nil, postgres.Error(err)
	}
	rows, err := r.db.Query(ctx, `
SELECT course.id, course.name, course.site_link,
       (SELECT count(*) FROM doctor_courses related WHERE related.course_id=course.id)
FROM doctor_courses completion
JOIN courses course ON course.id=completion.course_id
WHERE completion.doctor_id=$1
ORDER BY course.id`, doctorID)
	if err != nil {
		return nil, postgres.Error(err)
	}
	return postgres.Collect(rows, func(row pgx.Row) (*coursedto.Course, error) {
		var course coursedto.Course
		err := row.Scan(&course.ID, &course.Name, &course.SiteLink, &course.DoctorsCount)
		return &course, err
	})
}

func (r *Repository) AddDoctorCourse(ctx context.Context, doctorID, courseID int64) error {
	var doctorExists, courseExists bool
	if err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM doctors WHERE id=$1), EXISTS(SELECT 1 FROM courses WHERE id=$2)`, doctorID, courseID).Scan(&doctorExists, &courseExists); err != nil {
		return postgres.Error(err)
	}
	if !doctorExists || !courseExists {
		return domain_error.ErrNotFound
	}
	_, err := r.db.Exec(ctx, `INSERT INTO doctor_courses(doctor_id, course_id) VALUES($1, $2) ON CONFLICT (doctor_id, course_id) DO NOTHING`, doctorID, courseID)
	return postgres.Error(err)
}

func (r *Repository) DeleteDoctorCourse(ctx context.Context, doctorID, courseID int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM doctor_courses WHERE doctor_id=$1 AND course_id=$2`, doctorID, courseID)
	return postgres.DeleteResult(tag, err)
}
