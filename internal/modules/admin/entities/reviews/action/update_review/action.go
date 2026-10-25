package update_review

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/dto"
)

type Input = dto.UpdateInput

type Storage interface {
	UpdateReview(ctx context.Context, id int64, input dto.UpdateInput) (*dto.Review, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64, input Input) (*dto.Review, error) {
	return a.storage.UpdateReview(ctx, id, input)
}
