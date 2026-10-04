package create_course

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
)

type Action struct{ crud crud.Action }

func New(crudAction crud.Action) *Action { return &Action{crud: crudAction} }

func (a *Action) Execute(ctx context.Context, name, siteLink string) (map[string]any, error) {
	return a.crud.Create(ctx, "courses", map[string]any{
		"name":      name,
		"site_link": siteLink,
	})
}
