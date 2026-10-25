package get_city

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities/dto"
)

type Storage interface {
	GetCity(ctx context.Context, id int64) (*dto.City, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64) (*dto.City, error) {
	return a.storage.GetCity(ctx, id)
}
