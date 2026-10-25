package cities

import (
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities/action/create_city"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities/action/delete_city"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities/action/get_city"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities/action/list_cities"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities/action/update_city"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities/dal"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Aggregator struct {
	CreateCity *create_city.Action
	GetCity    *get_city.Action
	ListCities *list_cities.Action
	UpdateCity *update_city.Action
	DeleteCity *delete_city.Action
}

func New(db postgres.DB) *Aggregator {
	repository := dal.New(db)
	return &Aggregator{
		CreateCity: create_city.New(repository),
		GetCity:    get_city.New(repository),
		ListCities: list_cities.New(repository),
		UpdateCity: update_city.New(repository),
		DeleteCity: delete_city.New(repository),
	}
}
