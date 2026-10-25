package delete_review

import (
	"context"
)

type Storage interface {
	DeleteReview(ctx context.Context, id int64) error
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64) error {
	return a.storage.DeleteReview(ctx, id)
}
