package specialties

import (
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties/action/create_specialty"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties/action/delete_specialty"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties/action/get_specialty"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties/action/list_specialties"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties/action/update_specialty"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties/dal"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Aggregator struct {
	CreateSpecialty *create_specialty.Action
	GetSpecialty    *get_specialty.Action
	ListSpecialties *list_specialties.Action
	UpdateSpecialty *update_specialty.Action
	DeleteSpecialty *delete_specialty.Action
}

func New(db postgres.DB) *Aggregator {
	repository := dal.New(db)
	return &Aggregator{
		CreateSpecialty: create_specialty.New(repository),
		GetSpecialty:    get_specialty.New(repository),
		ListSpecialties: list_specialties.New(repository),
		UpdateSpecialty: update_specialty.New(repository),
		DeleteSpecialty: delete_specialty.New(repository),
	}
}
