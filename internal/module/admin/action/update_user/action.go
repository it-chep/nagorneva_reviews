package update_user

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/validation"
)

type Action struct{ crud crud.Action }

func New(crudAction crud.Action) *Action { return &Action{crud: crudAction} }

func (a *Action) Execute(ctx context.Context, id int64, email, rawPassword *string) (map[string]any, error) {
	values := map[string]any{}
	if email != nil {
		values["email"] = *email
	}
	if rawPassword != nil {
		passwordHash, err := validation.HashPassword(*rawPassword)
		if err != nil {
			return nil, err
		}
		values["password"] = passwordHash
	}
	return a.crud.Update(ctx, "users", id, values)
}
