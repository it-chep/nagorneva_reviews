package create_city

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities/dto"
)

type Storage interface {
	CreateCity(ctx context.Context, input dto.CreateInput) (*dto.City, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, name string, lat, lon float64) (*dto.City, error) {
	return a.storage.CreateCity(ctx, dto.CreateInput{Name: name, Lat: lat, Lon: lon})
}
