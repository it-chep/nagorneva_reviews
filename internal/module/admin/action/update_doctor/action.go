package update_doctor

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
)

type Input struct {
	Name                *string
	FullName            *string
	CityID              *int64
	SpecialtyID         *int64
	Lat                 *float64
	Lon                 *float64
	PersonalDataConsent *bool
	IsActive            *bool
}

type Action struct{ crud crud.Action }

func New(crudAction crud.Action) *Action { return &Action{crud: crudAction} }

func (a *Action) Execute(ctx context.Context, id int64, input Input) (map[string]any, error) {
	values := map[string]any{}
	if input.Name != nil {
		values["name"] = *input.Name
	}
	if input.FullName != nil {
		values["full_name"] = *input.FullName
	}
	if input.CityID != nil {
		values["city_id"] = *input.CityID
	}
	if input.SpecialtyID != nil {
		values["specialty_id"] = *input.SpecialtyID
	}
	if input.Lat != nil {
		values["lat"] = *input.Lat
	}
	if input.Lon != nil {
		values["lon"] = *input.Lon
	}
	if input.PersonalDataConsent != nil {
		values["personal_data_consent"] = *input.PersonalDataConsent
	}
	if input.IsActive != nil {
		values["is_active"] = *input.IsActive
	}
	return a.crud.Update(ctx, "doctors", id, values)
}
