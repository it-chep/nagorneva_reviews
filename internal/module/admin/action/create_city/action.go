package create_city

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
)

type Action struct{ crud crud.Action }

func New(crudAction crud.Action) *Action { return &Action{crud: crudAction} }

func (a *Action) Execute(ctx context.Context, name string, lat, lon float64) (map[string]any, error) {
	return a.crud.Create(ctx, "cities", map[string]any{"name": name, "lat": lat, "lon": lon})
}
