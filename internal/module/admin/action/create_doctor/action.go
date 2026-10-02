package create_doctor

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
)

type Input struct {
	Name                string
	FullName            string
	CityID              int64
	SpecialtyID         int64
	Lat                 float64
	Lon                 float64
	PersonalDataConsent bool
	IsActive            bool
}

type Action struct{ crud crud.Action }

func New(crudAction crud.Action) *Action { return &Action{crud: crudAction} }

func (a *Action) Execute(ctx context.Context, input Input) (map[string]any, error) {
	return a.crud.Create(ctx, "doctors", map[string]any{
		"name": input.Name, "full_name": input.FullName, "city_id": input.CityID, "specialty_id": input.SpecialtyID,
		"lat": input.Lat, "lon": input.Lon, "personal_data_consent": input.PersonalDataConsent, "is_active": input.IsActive,
	})
}
