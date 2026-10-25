package create_user

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/dto"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/service/password"
)

type Storage interface {
	CreateUser(ctx context.Context, input dto.CreateInput) (*dto.User, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, email, rawPassword string) (*dto.User, error) {
	hash, err := password.Hash(rawPassword)
	if err != nil {
		return nil, err
	}
	return a.storage.CreateUser(ctx, dto.CreateInput{Email: email, Password: hash})
}
