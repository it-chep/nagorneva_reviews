package filter_reviews

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/domain"
)

type Storage interface {
	Filter(ctx context.Context, courseID *int64, lat, lon *float64, radius float64) ([]domain.ReviewWithCourse, error)
}
type Action struct{ dal Storage }

func New(storage Storage) Action { return Action{dal: storage} }
func (a Action) Execute(ctx context.Context, courseID *int64, lat, lon *float64, radiusKm float64) ([]domain.ReviewWithCourse, error) {
	return a.dal.Filter(ctx, courseID, lat, lon, radiusKm)
}
