package list_reviews

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/dto"
)

type Storage interface {
	ListReviews(ctx context.Context) ([]dto.Review, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context) ([]dto.Review, error) {
	return a.storage.ListReviews(ctx)
}
