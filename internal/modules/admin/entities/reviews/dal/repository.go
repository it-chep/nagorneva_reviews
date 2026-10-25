package dal

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/dto"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Repository struct{ db postgres.DB }

func New(db postgres.DB) *Repository { return &Repository{db: db} }
func scan(row pgx.Row) (*dto.Review, error) {
	var v dto.Review
	err := row.Scan(&v.ID, &v.DoctorID, &v.Rating, &v.Comment, &v.CourseID, &v.IsActive, &v.CreatedAt)
	if err != nil {
		return nil, postgres.Error(err)
	}
	return &v, nil
}
func (r *Repository) GetReview(ctx context.Context, id int64) (*dto.Review, error) {
	return scan(r.db.QueryRow(ctx, `SELECT id, doctor_id, rating, comment, course_id, is_active, created_at FROM reviews WHERE id=$1`, id))
}
func (r *Repository) CreateReview(ctx context.Context, input dto.CreateInput) (*dto.Review, error) {
	return scan(r.db.QueryRow(ctx, `INSERT INTO reviews (doctor_id, rating, comment, course_id, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING id, doctor_id, rating, comment, course_id, is_active, created_at`, input.DoctorID, input.Rating, input.Comment, input.CourseID, input.IsActive))
}
func (r *Repository) UpdateReview(ctx context.Context, id int64, input dto.UpdateInput) (*dto.Review, error) {
	return scan(r.db.QueryRow(ctx, `UPDATE reviews SET rating=COALESCE($2, rating), comment=COALESCE($3, comment), course_id=COALESCE($4, course_id), is_active=COALESCE($5, is_active) WHERE id=$1 RETURNING id, doctor_id, rating, comment, course_id, is_active, created_at`, id, input.Rating, input.Comment, input.CourseID, input.IsActive))
}
func (r *Repository) ListReviews(ctx context.Context) ([]dto.Review, error) {
	rows, err := r.db.Query(ctx, `SELECT id, doctor_id, rating, comment, course_id, is_active, created_at FROM reviews ORDER BY id`)
	if err != nil {
		return nil, postgres.Error(err)
	}
	return postgres.Collect(rows, scan)
}
func (r *Repository) DeleteReview(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM reviews WHERE id=$1`, id)
	return postgres.DeleteResult(tag, err)
}
func (r *Repository) GetDoctorReviews(ctx context.Context, doctorID int64) ([]dto.Review, error) {
	rows, err := r.db.Query(ctx, `SELECT r.id,r.doctor_id,r.rating,r.comment,r.course_id,r.is_active,r.created_at,c.name FROM reviews r JOIN courses c ON c.id=r.course_id WHERE r.doctor_id=$1 ORDER BY r.created_at DESC`, doctorID)
	if err != nil {
		return nil, postgres.Error(err)
	}
	return postgres.Collect(rows, func(row pgx.Row) (*dto.Review, error) {
		var v dto.Review
		err := row.Scan(&v.ID, &v.DoctorID, &v.Rating, &v.Comment, &v.CourseID, &v.IsActive, &v.CreatedAt, &v.CourseName)
		return &v, err
	})
}
