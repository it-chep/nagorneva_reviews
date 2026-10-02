package update_city

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
)

type Action struct{ crud crud.Action }

func New(crudAction crud.Action) *Action { return &Action{crud: crudAction} }

func (a *Action) Execute(ctx context.Context, id int64, name *string, lat, lon *float64) (map[string]any, error) {
	values := map[string]any{}
	if name != nil {
		values["name"] = *name
	}
	if lat != nil {
		values["lat"] = *lat
	}
	if lon != nil {
		values["lon"] = *lon
	}
	return a.crud.Update(ctx, "cities", id, values)
}
