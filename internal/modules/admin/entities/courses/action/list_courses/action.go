package list_courses

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/dto"
)

type Storage interface {
	ListCourses(ctx context.Context) ([]dto.Course, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context) ([]dto.Course, error) {
	return a.storage.ListCourses(ctx)
}
