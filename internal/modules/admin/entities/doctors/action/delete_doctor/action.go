package delete_doctor

import (
	"context"
)

type Storage interface {
	DeleteDoctor(ctx context.Context, id int64) error
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64) error {
	return a.storage.DeleteDoctor(ctx, id)
}
