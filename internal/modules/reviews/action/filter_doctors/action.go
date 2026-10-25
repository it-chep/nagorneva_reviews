// Package filter_doctors contains the public doctor catalog filtering use case.
package filter_doctors

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/action/filter_doctors/dto"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/domain"
)

type Storage interface {
	Filter(ctx context.Context, cityIDs, specialtyIDs, courseIDs []int64, isActive, personalDataConsent *bool, reviewsSort ReviewsSort) ([]domain.FilteredDoctor, error)
}
type Action struct{ dal Storage }

type ReviewsSort = dto.ReviewsSort

const (
	ReviewsSortDefault    = dto.ReviewsSortDefault
	ReviewsSortDescending = dto.ReviewsSortDescending
	ReviewsSortAscending  = dto.ReviewsSortAscending
)

func New(storage Storage) Action { return Action{dal: storage} }

func (a Action) Execute(ctx context.Context, cityIDs, specialtyIDs, courseIDs []int64, isActive, personalDataConsent *bool, reviewsSort ReviewsSort) ([]domain.FilteredDoctor, error) {
	return a.dal.Filter(ctx, cityIDs, specialtyIDs, courseIDs, isActive, personalDataConsent, reviewsSort)
}
