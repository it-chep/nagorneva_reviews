package action_test

import (
	"context"
	"errors"
	"testing"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/action/add_doctor_course"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/action/delete_doctor_course"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/action/get_doctor_courses"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/domain_error"
)

func TestDoctorCoursesRejectInvalidIDsBeforeStorage(t *testing.T) {
	ctx := context.Background()
	for _, id := range []int64{0, -1} {
		if _, err := get_doctor_courses.New(nil).Execute(ctx, id); !errors.Is(err, domain_error.ErrInvalid) {
			t.Fatalf("GetDoctorCourses(%d): %v", id, err)
		}
	}
	for name, execute := range map[string]func(context.Context, int64, int64) error{
		"add":    add_doctor_course.New(nil).Execute,
		"delete": delete_doctor_course.New(nil).Execute,
	} {
		for _, ids := range [][2]int64{{0, 1}, {-1, 1}, {1, 0}, {1, -1}, {0, 0}} {
			if err := execute(ctx, ids[0], ids[1]); !errors.Is(err, domain_error.ErrInvalid) {
				t.Fatalf("%s(%d, %d): %v", name, ids[0], ids[1], err)
			}
		}
	}
}
