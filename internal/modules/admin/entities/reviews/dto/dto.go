package dto

import "time"

type Review struct {
	ID         int64
	DoctorID   int64
	Rating     int32
	Comment    string
	CourseID   int64
	IsActive   bool
	CreatedAt  time.Time
	CourseName string
}
type CreateInput struct {
	DoctorID int64
	Rating   int32
	Comment  string
	CourseID int64
	IsActive bool
}
type UpdateInput struct {
	Rating   *int32
	Comment  *string
	CourseID *int64
	IsActive *bool
}
