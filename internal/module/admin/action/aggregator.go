// Package action wires endpoint-specific admin actions into one module API.
package action

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/create_city"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/create_course"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/create_doctor"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/create_review"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/create_specialty"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/create_user"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/delete_city"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/delete_course"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/delete_doctor"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/delete_review"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/delete_specialty"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/delete_user"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/get_city"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/get_course"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/get_doctor"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/get_doctor_reviews"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/get_review"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/get_specialty"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/get_user"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/list_cities"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/list_courses"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/list_doctors"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/list_reviews"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/list_specialties"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/list_users"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/login"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/update_city"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/update_course"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/update_doctor"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/update_review"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/update_specialty"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/update_user"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/upload_doctor_photo"
	s3gateway "github.com/nagorneva/nagorneva_reviews/internal/module/admin/client/s3"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/auth"
)

// Aggregator exposes one action for every admin endpoint.
type Aggregator struct {
	CreateUser        *create_user.Action
	GetUser           *get_user.Action
	ListUsers         *list_users.Action
	UpdateUser        *update_user.Action
	DeleteUser        *delete_user.Action
	CreateCity        *create_city.Action
	GetCity           *get_city.Action
	ListCities        *list_cities.Action
	UpdateCity        *update_city.Action
	DeleteCity        *delete_city.Action
	CreateSpecialty   *create_specialty.Action
	GetSpecialty      *get_specialty.Action
	ListSpecialties   *list_specialties.Action
	UpdateSpecialty   *update_specialty.Action
	DeleteSpecialty   *delete_specialty.Action
	CreateCourse      *create_course.Action
	GetCourse         *get_course.Action
	ListCourses       *list_courses.Action
	UpdateCourse      *update_course.Action
	DeleteCourse      *delete_course.Action
	CreateDoctor      *create_doctor.Action
	GetDoctor         *get_doctor.Action
	ListDoctors       *list_doctors.Action
	UpdateDoctor      *update_doctor.Action
	DeleteDoctor      *delete_doctor.Action
	UploadDoctorPhoto *upload_doctor_photo.Action
	CreateReview      *create_review.Action
	GetReview         *get_review.Action
	ListReviews       *list_reviews.Action
	GetDoctorReviews  *get_doctor_reviews.Action
	UpdateReview      *update_review.Action
	DeleteReview      *delete_review.Action
	Login             *login.Action
}

func NewAggregator(db *pgxpool.Pool, token auth.Service, s3Client *s3gateway.Gateway) *Aggregator {
	crudAction := crud.New(db)
	return &Aggregator{
		CreateUser:        create_user.New(crudAction),
		GetUser:           get_user.New(crudAction),
		ListUsers:         list_users.New(crudAction),
		UpdateUser:        update_user.New(crudAction),
		DeleteUser:        delete_user.New(crudAction),
		CreateCity:        create_city.New(crudAction),
		GetCity:           get_city.New(crudAction),
		ListCities:        list_cities.New(crudAction),
		UpdateCity:        update_city.New(crudAction),
		DeleteCity:        delete_city.New(crudAction),
		CreateSpecialty:   create_specialty.New(crudAction),
		GetSpecialty:      get_specialty.New(crudAction),
		ListSpecialties:   list_specialties.New(crudAction),
		UpdateSpecialty:   update_specialty.New(crudAction),
		DeleteSpecialty:   delete_specialty.New(crudAction),
		CreateCourse:      create_course.New(crudAction),
		GetCourse:         get_course.New(crudAction),
		ListCourses:       list_courses.New(crudAction),
		UpdateCourse:      update_course.New(crudAction),
		DeleteCourse:      delete_course.New(crudAction),
		CreateDoctor:      create_doctor.New(crudAction),
		GetDoctor:         get_doctor.New(crudAction),
		ListDoctors:       list_doctors.New(crudAction),
		UpdateDoctor:      update_doctor.New(crudAction),
		DeleteDoctor:      delete_doctor.New(crudAction),
		UploadDoctorPhoto: upload_doctor_photo.New(crudAction, s3Client),
		CreateReview:      create_review.New(crudAction),
		GetReview:         get_review.New(crudAction),
		ListReviews:       list_reviews.New(crudAction),
		GetDoctorReviews:  get_doctor_reviews.New(crudAction),
		UpdateReview:      update_review.New(crudAction),
		DeleteReview:      delete_review.New(crudAction),
		Login:             login.New(db, token),
	}
}
