package delete_city

import (
	"context"
)

type Storage interface {
	DeleteCity(ctx context.Context, id int64) error
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64) error {
	return a.storage.DeleteCity(ctx, id)
}
