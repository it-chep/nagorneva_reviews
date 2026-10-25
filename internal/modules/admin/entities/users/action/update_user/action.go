package update_user

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/dto"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/service/password"
)

type Storage interface {
	UpdateUser(ctx context.Context, id int64, input dto.UpdateInput) (*dto.User, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64, email, rawPassword *string) (*dto.User, error) {
	var hash *string
	if rawPassword != nil {
		value, err := password.Hash(*rawPassword)
		if err != nil {
			return nil, err
		}
		hash = &value
	}
	return a.storage.UpdateUser(ctx, id, dto.UpdateInput{Email: email, Password: hash})
}
