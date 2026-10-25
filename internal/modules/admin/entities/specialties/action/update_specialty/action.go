package update_specialty

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties/dto"
)

type Storage interface {
	UpdateSpecialty(ctx context.Context, id int64, input dto.UpdateInput) (*dto.Specialty, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64, name *string) (*dto.Specialty, error) {
	return a.storage.UpdateSpecialty(ctx, id, dto.UpdateInput{Name: name})
}
