package adminv1

import (
	citydto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities/dto"
	coursedto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/dto"
	doctordto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/dto"
	reviewdto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/dto"
	specialtydto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties/dto"
	userdto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/dto"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func user(v *userdto.User) *adminpb.User {
	if v == nil {
		return nil
	}
	return &adminpb.User{Id: v.ID, Email: v.Email}
}

func city(v *citydto.City) *adminpb.City {
	if v == nil {
		return nil
	}
	return &adminpb.City{Id: v.ID, Name: v.Name, Lat: v.Lat, Lon: v.Lon, DoctorsCount: v.DoctorsCount}
}

func specialty(v *specialtydto.Specialty) *adminpb.Specialty {
	if v == nil {
		return nil
	}
	return &adminpb.Specialty{Id: v.ID, Name: v.Name, DoctorsCount: v.DoctorsCount}
}

func course(v *coursedto.Course) *adminpb.Course {
	if v == nil {
		return nil
	}
	return &adminpb.Course{
		Id:           v.ID,
		Name:         v.Name,
		DoctorsCount: v.DoctorsCount,
		SiteLink:     v.SiteLink,
	}
}

func doctor(v *doctordto.Doctor) *adminpb.Doctor {
	if v == nil {
		return nil
	}
	return &adminpb.Doctor{
		Id:                  v.ID,
		Name:                v.Name,
		FullName:            v.FullName,
		Photo:               v.Photo,
		PersonalDataConsent: v.PersonalDataConsent,
		City: &adminpb.City{
			Id:   v.CityID,
			Name: v.CityName,
			Lat:  v.CityLat,
			Lon:  v.CityLon,
		},
		Specialty: &adminpb.Specialty{
			Id:   v.SpecialtyID,
			Name: v.SpecialtyName,
		},
		Lat:          v.Lat,
		Lon:          v.Lon,
		IsActive:     v.IsActive,
		ReviewsCount: v.ReviewsCount,
	}
}

func review(v *reviewdto.Review) *adminpb.Review {
	if v == nil {
		return nil
	}
	out := &adminpb.Review{Id: v.ID, DoctorId: v.DoctorID, Rating: int32(v.Rating), Comment: v.Comment, CourseId: v.CourseID, IsActive: v.IsActive, CourseName: v.CourseName}
	if !v.CreatedAt.IsZero() {
		out.CreatedAt = timestamppb.New(v.CreatedAt)
	}
	return out
}

func users(items []userdto.User) []*adminpb.User {
	out := make([]*adminpb.User, 0, len(items))
	for _, item := range items {
		out = append(out, user(&item))
	}
	return out
}

func cities(items []citydto.City) []*adminpb.City {
	out := make([]*adminpb.City, 0, len(items))
	for _, item := range items {
		out = append(out, city(&item))
	}
	return out
}

func specialties(items []specialtydto.Specialty) []*adminpb.Specialty {
	out := make([]*adminpb.Specialty, 0, len(items))
	for _, item := range items {
		out = append(out, specialty(&item))
	}
	return out
}

func courses(items []coursedto.Course) []*adminpb.Course {
	out := make([]*adminpb.Course, 0, len(items))
	for _, item := range items {
		out = append(out, course(&item))
	}
	return out
}

func doctors(items []doctordto.Doctor) []*adminpb.Doctor {
	out := make([]*adminpb.Doctor, 0, len(items))
	for _, item := range items {
		out = append(out, doctor(&item))
	}
	return out
}

func reviews(items []reviewdto.Review) []*adminpb.Review {
	out := make([]*adminpb.Review, 0, len(items))
	for _, item := range items {
		out = append(out, review(&item))
	}
	return out
}
