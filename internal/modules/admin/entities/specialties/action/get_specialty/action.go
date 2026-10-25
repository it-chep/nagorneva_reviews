package get_specialty

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties/dto"
)

type Storage interface {
	GetSpecialty(ctx context.Context, id int64) (*dto.Specialty, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64) (*dto.Specialty, error) {
	return a.storage.GetSpecialty(ctx, id)
}
