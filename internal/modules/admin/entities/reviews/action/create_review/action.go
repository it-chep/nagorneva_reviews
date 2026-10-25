package create_review

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/dto"
)

type Storage interface {
	CreateReview(ctx context.Context, input dto.CreateInput) (*dto.Review, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, doctorID int64, rating int32, comment string, courseID int64, isActive bool) (*dto.Review, error) {
	return a.storage.CreateReview(ctx, dto.CreateInput{DoctorID: doctorID, Rating: rating, Comment: comment, CourseID: courseID, IsActive: isActive})
}
