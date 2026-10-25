package dal

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/dto"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Repository struct{ db postgres.DB }

func New(db postgres.DB) *Repository { return &Repository{db: db} }
func scan(row pgx.Row) (*dto.Course, error) {
	var v dto.Course
	err := row.Scan(&v.ID, &v.Name, &v.SiteLink)
	if err != nil {
		return nil, postgres.Error(err)
	}
	return &v, nil
}
func (r *Repository) GetCourse(ctx context.Context, id int64) (*dto.Course, error) {
	return scan(r.db.QueryRow(ctx, `SELECT id, name, site_link FROM courses WHERE id=$1`, id))
}
func (r *Repository) CreateCourse(ctx context.Context, input dto.CreateInput) (*dto.Course, error) {
	return scan(r.db.QueryRow(ctx, `INSERT INTO courses (name, site_link) VALUES ($1, $2) RETURNING id, name, site_link`, input.Name, input.SiteLink))
}
func (r *Repository) UpdateCourse(ctx context.Context, id int64, input dto.UpdateInput) (*dto.Course, error) {
	return scan(r.db.QueryRow(ctx, `UPDATE courses SET name=COALESCE($2, name), site_link=COALESCE($3, site_link) WHERE id=$1 RETURNING id, name, site_link`, id, input.Name, input.SiteLink))
}
func (r *Repository) ListCourses(ctx context.Context) ([]dto.Course, error) {
	rows, err := r.db.Query(ctx, `SELECT entity.id, entity.name, entity.site_link, (SELECT count(*) FROM doctor_courses related WHERE related.course_id = entity.id) FROM courses entity ORDER BY entity.id`)
	if err != nil {
		return nil, postgres.Error(err)
	}
	return postgres.Collect(rows, func(row pgx.Row) (*dto.Course, error) {
		var v dto.Course
		err := row.Scan(&v.ID, &v.Name, &v.SiteLink, &v.DoctorsCount)
		return &v, err
	})
}
func (r *Repository) DeleteCourse(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM courses WHERE id=$1`, id)
	return postgres.DeleteResult(tag, err)
}
