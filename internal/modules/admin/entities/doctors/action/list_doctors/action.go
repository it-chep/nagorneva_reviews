package list_doctors

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/dto"
)

type Storage interface {
	ListDoctors(ctx context.Context) ([]dto.Doctor, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context) ([]dto.Doctor, error) {
	return a.storage.ListDoctors(ctx)
}
