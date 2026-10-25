package list_cities

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities/dto"
)

type Storage interface {
	ListCities(ctx context.Context) ([]dto.City, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context) ([]dto.City, error) {
	return a.storage.ListCities(ctx)
}
