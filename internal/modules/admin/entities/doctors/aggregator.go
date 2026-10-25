package doctors

import (
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/action/create_doctor"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/action/delete_doctor"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/action/get_doctor"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/action/list_doctors"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/action/update_doctor"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/action/upload_doctor_photo"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/dal"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Aggregator struct {
	CreateDoctor      *create_doctor.Action
	GetDoctor         *get_doctor.Action
	ListDoctors       *list_doctors.Action
	UpdateDoctor      *update_doctor.Action
	DeleteDoctor      *delete_doctor.Action
	UploadDoctorPhoto *upload_doctor_photo.Action
}

func New(db postgres.DB, images upload_doctor_photo.ImageStorage) *Aggregator {
	repository := dal.New(db)
	return &Aggregator{
		CreateDoctor:      create_doctor.New(repository),
		GetDoctor:         get_doctor.New(repository),
		ListDoctors:       list_doctors.New(repository),
		UpdateDoctor:      update_doctor.New(repository),
		DeleteDoctor:      delete_doctor.New(repository),
		UploadDoctorPhoto: upload_doctor_photo.New(repository, images),
	}
}
