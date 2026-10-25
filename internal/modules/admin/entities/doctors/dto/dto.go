package dto

type Doctor struct {
	ID                  int64
	Name                string
	FullName            string
	CityID              int64
	SpecialtyID         int64
	Lat                 float64
	Lon                 float64
	PersonalDataConsent bool
	IsActive            bool
	Photo               string
	CityName            string
	CityLat             float64
	CityLon             float64
	SpecialtyName       string
	ReviewsCount        int64
}
type CreateInput struct {
	Name                string
	FullName            string
	CityID              int64
	SpecialtyID         int64
	Lat                 float64
	Lon                 float64
	PersonalDataConsent bool
	IsActive            bool
}
type UpdateInput struct {
	Name                *string
	FullName            *string
	CityID              *int64
	SpecialtyID         *int64
	Lat                 *float64
	Lon                 *float64
	PersonalDataConsent *bool
	IsActive            *bool
}
