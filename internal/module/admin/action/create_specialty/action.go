package create_specialty

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
)

type Action struct{ crud crud.Action }

func New(crudAction crud.Action) *Action { return &Action{crud: crudAction} }

func (a *Action) Execute(ctx context.Context, name string) (map[string]any, error) {
	return a.crud.Create(ctx, "specialties", map[string]any{"name": name})
}
