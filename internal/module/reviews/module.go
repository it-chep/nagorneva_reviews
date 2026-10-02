// Package reviews composes actions for public reviews and map endpoints.
package reviews

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nagorneva/nagorneva_reviews/internal/module/reviews/action"
)

// Module is the entry point for public reviews use cases.
type Module struct {
	Actions *action.Aggregator
}

// New constructs all public reviews actions.
func New(db *pgxpool.Pool) *Module {
	return &Module{Actions: action.NewAggregator(db)}
}
