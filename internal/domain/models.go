package domain

import "time"

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
type City struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}
type Specialty struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
type Course struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
type Doctor struct {
	ID                  int64   `json:"id"`
	Name                string  `json:"name"`
	FullName            string  `json:"full_name"`
	Photo               string  `json:"photo"`
	PersonalDataConsent bool    `json:"personal_data_consent"`
	CityID              int64   `json:"city_id"`
	SpecialtyID         int64   `json:"specialty_id"`
	Lat                 float64 `json:"lat"`
	Lon                 float64 `json:"lon"`
	IsActive            bool    `json:"is_active"`
}
type Review struct {
	ID        int64     `json:"id"`
	DoctorID  int64     `json:"doctor_id"`
	Rating    int16     `json:"rating"`
	Comment   string    `json:"comment"`
	CourseID  int64     `json:"course_id"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}
type ReviewWithCourse struct {
	Review
	CourseName string `json:"course_name"`
}
type DoctorOnMap struct {
	Doctor
	City    string             `json:"city"`
	Reviews []ReviewWithCourse `json:"reviews"`
}
