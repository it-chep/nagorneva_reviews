package list_users

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/dto"
)

type Storage interface {
	ListUsers(ctx context.Context) ([]dto.User, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context) ([]dto.User, error) {
	return a.storage.ListUsers(ctx)
}
