package get_doctor_courses

import (
	"context"

	coursedto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/dto"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/domain_error"
)

type Storage interface {
	GetDoctorCourses(context.Context, int64) ([]coursedto.Course, error)
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }

func (a *Action) Execute(ctx context.Context, doctorID int64) ([]coursedto.Course, error) {
	if doctorID <= 0 {
		return nil, domain_error.ErrInvalid
	}
	return a.storage.GetDoctorCourses(ctx, doctorID)
}
