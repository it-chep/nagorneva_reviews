package create_review

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
)

type Action struct{ crud crud.Action }

func New(crudAction crud.Action) *Action { return &Action{crud: crudAction} }

func (a *Action) Execute(ctx context.Context, doctorID int64, rating int32, comment string, courseID int64, isActive bool) (map[string]any, error) {
	return a.crud.Create(ctx, "reviews", map[string]any{
		"doctor_id": doctorID, "rating": rating, "comment": comment, "course_id": courseID, "is_active": isActive,
	})
}
