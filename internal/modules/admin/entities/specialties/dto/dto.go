package dto

type Specialty struct {
	ID           int64
	Name         string
	DoctorsCount int64
}
type CreateInput struct {
	Name string
}
type UpdateInput struct {
	Name *string
}
