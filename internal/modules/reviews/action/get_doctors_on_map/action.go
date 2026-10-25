package get_doctors_on_map

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/domain"
)

type Storage interface {
	Get(ctx context.Context, lat, lon, radius float64) ([]domain.DoctorOnMap, error)
}
type Action struct{ dal Storage }

func New(storage Storage) Action { return Action{dal: storage} }
func (a Action) Execute(ctx context.Context, lat, lon, radiusKm float64) ([]domain.DoctorOnMap, error) {
	return a.dal.Get(ctx, lat, lon, radiusKm)
}
