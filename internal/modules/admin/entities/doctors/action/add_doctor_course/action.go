package add_doctor_course

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/pkg/domain_error"
)

type Storage interface {
	AddDoctorCourse(context.Context, int64, int64) error
}
type Action struct{ storage Storage }

func New(storage Storage) *Action { return &Action{storage: storage} }

func (a *Action) Execute(ctx context.Context, doctorID, courseID int64) error {
	if doctorID <= 0 || courseID <= 0 {
		return domain_error.ErrInvalid
	}
	return a.storage.AddDoctorCourse(ctx, doctorID, courseID)
}
