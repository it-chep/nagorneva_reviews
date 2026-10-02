package update_review

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
)

type Input struct {
	Rating   *int32
	Comment  *string
	CourseID *int64
	IsActive *bool
}

type Action struct{ crud crud.Action }

func New(crudAction crud.Action) *Action { return &Action{crud: crudAction} }

func (a *Action) Execute(ctx context.Context, id int64, input Input) (map[string]any, error) {
	values := map[string]any{}
	if input.Rating != nil {
		values["rating"] = *input.Rating
	}
	if input.Comment != nil {
		values["comment"] = *input.Comment
	}
	if input.CourseID != nil {
		values["course_id"] = *input.CourseID
	}
	if input.IsActive != nil {
		values["is_active"] = *input.IsActive
	}
	return a.crud.Update(ctx, "reviews", id, values)
}
