package update_city

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities/dto"
)

type Storage interface {
	UpdateCity(ctx context.Context, id int64, input dto.UpdateInput) (*dto.City, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64, name *string, lat, lon *float64) (*dto.City, error) {
	return a.storage.UpdateCity(ctx, id, dto.UpdateInput{Name: name, Lat: lat, Lon: lon})
}
