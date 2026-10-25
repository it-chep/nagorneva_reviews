package get_doctor

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/dto"
)

type Storage interface {
	GetDoctor(ctx context.Context, id int64) (*dto.Doctor, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64) (*dto.Doctor, error) {
	return a.storage.GetDoctor(ctx, id)
}
