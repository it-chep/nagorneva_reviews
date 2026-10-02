package create_user

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/validation"
)

type Action struct{ crud crud.Action }

func New(crudAction crud.Action) *Action { return &Action{crud: crudAction} }

func (a *Action) Execute(ctx context.Context, email, rawPassword string) (map[string]any, error) {
	passwordHash, err := validation.HashPassword(rawPassword)
	if err != nil {
		return nil, err
	}
	return a.crud.Create(ctx, "users", map[string]any{"email": email, "password": passwordHash})
}
