package courses

import (
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/action/create_course"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/action/delete_course"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/action/get_course"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/action/list_courses"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/action/update_course"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/dal"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Aggregator struct {
	CreateCourse *create_course.Action
	GetCourse    *get_course.Action
	ListCourses  *list_courses.Action
	UpdateCourse *update_course.Action
	DeleteCourse *delete_course.Action
}

func New(db postgres.DB) *Aggregator {
	repository := dal.New(db)
	return &Aggregator{
		CreateCourse: create_course.New(repository),
		GetCourse:    get_course.New(repository),
		ListCourses:  list_courses.New(repository),
		UpdateCourse: update_course.New(repository),
		DeleteCourse: delete_course.New(repository),
	}
}
