package create_specialty

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties/dto"
)

type Storage interface {
	CreateSpecialty(ctx context.Context, input dto.CreateInput) (*dto.Specialty, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, name string) (*dto.Specialty, error) {
	return a.storage.CreateSpecialty(ctx, dto.CreateInput{Name: name})
}
