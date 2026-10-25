// Package reviews composes actions for public reviews and map endpoints.
package reviews

import (
	"github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/action"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

// Module is the entry point for public reviews use cases.
type Module struct {
	Actions *action.Aggregator
}

// New constructs all public reviews actions.
func New(db postgres.DB) *Module {
	return &Module{Actions: action.NewAggregator(db)}
}
