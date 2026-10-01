package filter_reviews

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nagorneva/nagorneva_reviews/internal/domain"
)

type Action struct{ dal DAL }

func New(db *pgxpool.Pool) Action { return Action{DAL{db: db}} }
func (a Action) Execute(ctx context.Context, courseID *int64, lat, lon *float64, radiusKm float64) ([]domain.ReviewWithCourse, error) {
	return a.dal.Filter(ctx, courseID, lat, lon, radiusKm)
}
