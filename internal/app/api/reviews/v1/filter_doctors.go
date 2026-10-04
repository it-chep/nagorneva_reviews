package reviewsv1

import (
	"context"

	reviewspb "github.com/nagorneva/nagorneva_reviews/gen/go/reviews"
	"github.com/nagorneva/nagorneva_reviews/internal/domain"
	filterdoctors "github.com/nagorneva/nagorneva_reviews/internal/module/reviews/action/filter_doctors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) FilterDoctors(ctx context.Context, req *reviewspb.FilterDoctorsRequest) (*reviewspb.FilterDoctorsResponse, error) {
	if err := validateFilterIDs("city_ids", req.GetCityIds()); err != nil {
		return nil, err
	}
	if err := validateFilterIDs("specialty_ids", req.GetSpecialtyIds()); err != nil {
		return nil, err
	}
	if err := validateFilterIDs("course_ids", req.GetCourseIds()); err != nil {
		return nil, err
	}

	reviewsSort, err := reviewSort(req.GetReviewsSort())
	if err != nil {
		return nil, err
	}
	doctors, err := s.module.Actions.FilterDoctors.Execute(ctx, req.GetCityIds(), req.GetSpecialtyIds(), req.GetCourseIds(), req.IsActive, req.PersonalDataConsent, reviewsSort)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	result := make([]*reviewspb.FilteredDoctor, 0, len(doctors))
	for _, doctor := range doctors {
		result = append(result, filteredDoctor(doctor))
	}
	return &reviewspb.FilterDoctorsResponse{Doctors: result, DoctorCount: int64(len(result))}, nil
}

func reviewSort(value reviewspb.ReviewsSort) (filterdoctors.ReviewsSort, error) {
	switch value {
	case reviewspb.ReviewsSort_REVIEWS_SORT_UNSPECIFIED:
		return filterdoctors.ReviewsSortDefault, nil
	case reviewspb.ReviewsSort_REVIEWS_SORT_DESC:
		return filterdoctors.ReviewsSortDescending, nil
	case reviewspb.ReviewsSort_REVIEWS_SORT_ASC:
		return filterdoctors.ReviewsSortAscending, nil
	default:
		return filterdoctors.ReviewsSortDefault, status.Error(codes.InvalidArgument, "reviews_sort has an unsupported value")
	}
}

func validateFilterIDs(field string, ids []int64) error {
	for _, id := range ids {
		if id <= 0 {
			return status.Errorf(codes.InvalidArgument, "%s must contain only positive IDs", field)
		}
	}
	return nil
}

func filteredDoctor(doctor domain.FilteredDoctor) *reviewspb.FilteredDoctor {
	completedCourses := make([]*reviewspb.CompletedCourse, 0, len(doctor.CompletedCourses))
	for _, course := range doctor.CompletedCourses {
		completedCourses = append(completedCourses, &reviewspb.CompletedCourse{Id: course.ID, Name: course.Name})
	}
	return &reviewspb.FilteredDoctor{
		Id:    doctor.ID,
		Name:  doctor.Name,
		Photo: doctor.Photo,
		City: &reviewspb.FilterDoctorCity{
			Id: doctor.City.ID, Name: doctor.City.Name, Lat: doctor.City.Lat, Lon: doctor.City.Lon,
		},
		Specialty: &reviewspb.FilterDoctorSpecialty{
			Id: doctor.Specialty.ID, Name: doctor.Specialty.Name,
		},
		CompletedCourses:    completedCourses,
		IsActive:            doctor.IsActive,
		PersonalDataConsent: doctor.PersonalDataConsent,
		ReviewsCount:        doctor.ReviewsCount,
	}
}
