package entities

import (
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/action/login"
	logindal "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/action/login/dal"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/action/upload_doctor_photo"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Aggregator struct {
	Users       *users.Aggregator
	Cities      *cities.Aggregator
	Specialties *specialties.Aggregator
	Courses     *courses.Aggregator
	Doctors     *doctors.Aggregator
	Reviews     *reviews.Aggregator
	Login       *login.Action
}

func New(db postgres.DB, token login.TokenIssuer, images upload_doctor_photo.ImageStorage) *Aggregator {
	return &Aggregator{
		Users:       users.New(db),
		Cities:      cities.New(db),
		Specialties: specialties.New(db),
		Courses:     courses.New(db),
		Doctors:     doctors.New(db, images),
		Reviews:     reviews.New(db),
		Login:       login.New(logindal.New(db), token),
	}
}
