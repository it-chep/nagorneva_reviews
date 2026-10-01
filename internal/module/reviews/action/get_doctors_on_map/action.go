package get_doctors_on_map

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nagorneva/nagorneva_reviews/internal/domain"
)

type Action struct{ dal DAL }

func New(db *pgxpool.Pool) Action { return Action{dal: DAL{db: db}} }
func (a Action) Execute(ctx context.Context, lat, lon, radiusKm float64) ([]domain.DoctorOnMap, error) {
	return a.dal.Get(ctx, lat, lon, radiusKm)
}
