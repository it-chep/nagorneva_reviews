package get_doctor_reviews

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/dto"
)

type Storage interface {
	GetDoctorReviews(context.Context, int64) ([]dto.Review, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, doctorID int64) ([]dto.Review, error) {
	return a.storage.GetDoctorReviews(ctx, doctorID)
}
