package get_user

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/dto"
)

type Storage interface {
	GetUser(ctx context.Context, id int64) (*dto.User, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64) (*dto.User, error) {
	return a.storage.GetUser(ctx, id)
}
