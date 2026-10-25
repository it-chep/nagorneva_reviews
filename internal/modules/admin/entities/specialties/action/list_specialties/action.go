package list_specialties

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties/dto"
)

type Storage interface {
	ListSpecialties(ctx context.Context) ([]dto.Specialty, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context) ([]dto.Specialty, error) {
	return a.storage.ListSpecialties(ctx)
}
