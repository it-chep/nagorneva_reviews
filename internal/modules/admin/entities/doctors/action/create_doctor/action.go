package create_doctor

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/dto"
)

type Input = dto.CreateInput

type Storage interface {
	CreateDoctor(ctx context.Context, input dto.CreateInput) (*dto.Doctor, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, input Input) (*dto.Doctor, error) {
	return a.storage.CreateDoctor(ctx, input)
}
