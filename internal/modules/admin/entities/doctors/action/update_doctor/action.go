package update_doctor

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/dto"
)

type Input = dto.UpdateInput

type Storage interface {
	UpdateDoctor(ctx context.Context, id int64, input dto.UpdateInput) (*dto.Doctor, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64, input Input) (*dto.Doctor, error) {
	return a.storage.UpdateDoctor(ctx, id, input)
}
