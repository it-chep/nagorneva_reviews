package dto

type City struct {
	ID           int64
	Name         string
	Lat          float64
	Lon          float64
	DoctorsCount int64
}
type CreateInput struct {
	Name string
	Lat  float64
	Lon  float64
}
type UpdateInput struct {
	Name *string
	Lat  *float64
	Lon  *float64
}
