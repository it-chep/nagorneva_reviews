package create_course

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/dto"
)

type Storage interface {
	CreateCourse(ctx context.Context, input dto.CreateInput) (*dto.Course, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, name, siteLink string) (*dto.Course, error) {
	return a.storage.CreateCourse(ctx, dto.CreateInput{Name: name, SiteLink: siteLink})
}
