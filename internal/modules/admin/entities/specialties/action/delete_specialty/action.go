package delete_specialty

import (
	"context"
)

type Storage interface {
	DeleteSpecialty(ctx context.Context, id int64) error
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64) error {
	return a.storage.DeleteSpecialty(ctx, id)
}
