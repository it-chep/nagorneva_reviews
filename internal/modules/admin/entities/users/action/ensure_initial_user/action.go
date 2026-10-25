package ensure_initial_user

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/pkg/password"
)

type Storage interface {
	ExistsByEmail(context.Context, string) (bool, error)
	CreateInitialUser(context.Context, string, string) error
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, email, rawPassword string) error {
	if email == "" || rawPassword == "" {
		return nil
	}
	exists, err := a.storage.ExistsByEmail(ctx, email)
	if err != nil || exists {
		return err
	}
	hash, err := password.Hash(rawPassword)
	if err != nil {
		return err
	}
	return a.storage.CreateInitialUser(ctx, email, hash)
}
