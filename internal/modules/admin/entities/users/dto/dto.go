package dto

type User struct {
	ID    int64
	Email string
}
type CreateInput struct {
	Email    string
	Password string
}
type UpdateInput struct {
	Email    *string
	Password *string
}
