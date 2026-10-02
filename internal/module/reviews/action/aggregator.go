// Package action wires public reviews endpoint actions into the module.
package action

import (
	"github.com/jackc/pgx/v5/pgxpool"
	filterreviews "github.com/nagorneva/nagorneva_reviews/internal/module/reviews/action/filter_reviews"
	getdoctorsonmap "github.com/nagorneva/nagorneva_reviews/internal/module/reviews/action/get_doctors_on_map"
)

// Aggregator exposes one action per public reviews endpoint.
type Aggregator struct {
	FilterReviews   filterreviews.Action
	GetDoctorsOnMap getdoctorsonmap.Action
}

func NewAggregator(db *pgxpool.Pool) *Aggregator {
	return &Aggregator{
		FilterReviews:   filterreviews.New(db),
		GetDoctorsOnMap: getdoctorsonmap.New(db),
	}
}
