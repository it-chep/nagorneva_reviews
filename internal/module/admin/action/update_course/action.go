package update_course

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
)

type Action struct{ crud crud.Action }

func New(crudAction crud.Action) *Action { return &Action{crud: crudAction} }

func (a *Action) Execute(ctx context.Context, id int64, name, siteLink *string) (map[string]any, error) {
	values := map[string]any{}
	if name != nil {
		values["name"] = *name
	}
	if siteLink != nil {
		values["site_link"] = *siteLink
	}
	return a.crud.Update(ctx, "courses", id, values)
}
