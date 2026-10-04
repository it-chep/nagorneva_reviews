// Package filter_doctors contains the public doctor catalog filtering use case.
package filter_doctors

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nagorneva/nagorneva_reviews/internal/domain"
)

type Action struct{ dal DAL }

type ReviewsSort int32

const (
	ReviewsSortDefault ReviewsSort = iota
	ReviewsSortDescending
	ReviewsSortAscending
)

func New(db *pgxpool.Pool) Action { return Action{dal: DAL{db: db}} }

func (a Action) Execute(ctx context.Context, cityIDs, specialtyIDs, courseIDs []int64, isActive, personalDataConsent *bool, reviewsSort ReviewsSort) ([]domain.FilteredDoctor, error) {
	return a.dal.Filter(ctx, cityIDs, specialtyIDs, courseIDs, isActive, personalDataConsent, reviewsSort)
}
