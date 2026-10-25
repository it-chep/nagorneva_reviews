package update_course

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/dto"
)

type Storage interface {
	UpdateCourse(ctx context.Context, id int64, input dto.UpdateInput) (*dto.Course, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }
func (a *Action) Execute(ctx context.Context, id int64, name, siteLink *string) (*dto.Course, error) {
	return a.storage.UpdateCourse(ctx, id, dto.UpdateInput{Name: name, SiteLink: siteLink})
}
