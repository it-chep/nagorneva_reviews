package get_review

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/dto"
)

type Storage interface {
	GetReview(ctx context.Context, id int64) (*dto.Review, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64) (*dto.Review, error) {
	return a.storage.GetReview(ctx, id)
}
