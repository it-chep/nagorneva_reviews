// Package action wires public reviews endpoint actions into the module.
package action

import (
	filterdoctors "github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/action/filter_doctors"
	filterdoctorsdal "github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/action/filter_doctors/dal"
	filterreviews "github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/action/filter_reviews"
	filterreviewsdal "github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/action/filter_reviews/dal"
	getdoctorsonmap "github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/action/get_doctors_on_map"
	getdoctorsonmapdal "github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/action/get_doctors_on_map/dal"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

// Aggregator exposes one action per public reviews endpoint.
type Aggregator struct {
	FilterDoctors   filterdoctors.Action
	FilterReviews   filterreviews.Action
	GetDoctorsOnMap getdoctorsonmap.Action
}

func NewAggregator(db postgres.DB) *Aggregator {
	return &Aggregator{
		FilterDoctors:   filterdoctors.New(filterdoctorsdal.New(db)),
		FilterReviews:   filterreviews.New(filterreviewsdal.New(db)),
		GetDoctorsOnMap: getdoctorsonmap.New(getdoctorsonmapdal.New(db)),
	}
}
