package get_course

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/dto"
)

type Storage interface {
	GetCourse(ctx context.Context, id int64) (*dto.Course, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64) (*dto.Course, error) {
	return a.storage.GetCourse(ctx, id)
}
