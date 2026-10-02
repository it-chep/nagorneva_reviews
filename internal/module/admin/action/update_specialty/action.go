package update_specialty

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
)

type Action struct{ crud crud.Action }

func New(crudAction crud.Action) *Action { return &Action{crud: crudAction} }

func (a *Action) Execute(ctx context.Context, id int64, name *string) (map[string]any, error) {
	values := map[string]any{}
	if name != nil {
		values["name"] = *name
	}
	return a.crud.Update(ctx, "specialties", id, values)
}
