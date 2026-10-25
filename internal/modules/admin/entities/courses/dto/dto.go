package dto

type Course struct {
	ID           int64
	Name         string
	SiteLink     string
	DoctorsCount int64
}
type CreateInput struct {
	Name     string
	SiteLink string
}
type UpdateInput struct {
	Name     *string
	SiteLink *string
}
