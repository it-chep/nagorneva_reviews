package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
	filterdoctors "github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/action/filter_doctors"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) FilterDoctors(ctx context.Context, req *adminpb.FilterDoctorsRequest) (*adminpb.FilterDoctorsResponse, error) {
	for field, ids := range map[string][]int64{
		"city_ids":      req.GetCityIds(),
		"specialty_ids": req.GetSpecialtyIds(),
		"course_ids":    req.GetCourseIds(),
	} {
		for _, id := range ids {
			if id <= 0 {
				return nil, status.Errorf(codes.InvalidArgument, "%s must contain only positive IDs", field)
			}
		}
	}
	sort := req.GetReviewsSort()
	if sort < adminpb.ReviewsSort_REVIEWS_SORT_UNSPECIFIED || sort > adminpb.ReviewsSort_REVIEWS_SORT_ASC {
		return nil, status.Error(codes.InvalidArgument, "reviews_sort has an unsupported value")
	}
	items, err := s.module.Actions.Doctors.FilterDoctors.Execute(ctx, req.GetCityIds(), req.GetSpecialtyIds(), req.GetCourseIds(), req.IsActive, req.PersonalDataConsent, filterdoctors.ReviewsSort(sort))
	if err != nil {
		return nil, rpcError(err)
	}
	result := make([]*adminpb.FilteredDoctor, 0, len(items))
	for _, item := range items {
		result = append(result, filteredDoctor(item))
	}
	return &adminpb.FilterDoctorsResponse{Doctors: result, DoctorCount: int64(len(result))}, nil
}

func filteredDoctor(item domain.FilteredDoctor) *adminpb.FilteredDoctor {
	completed := make([]*adminpb.Course, 0, len(item.CompletedCourses))
	for _, course := range item.CompletedCourses {
		completed = append(completed, &adminpb.Course{Id: course.ID, Name: course.Name, SiteLink: course.SiteLink})
	}
	return &adminpb.FilteredDoctor{
		Id: item.ID, Name: item.Name, Photo: item.Photo,
		City:             &adminpb.City{Id: item.City.ID, Name: item.City.Name, Lat: item.City.Lat, Lon: item.City.Lon},
		Specialty:        &adminpb.Specialty{Id: item.Specialty.ID, Name: item.Specialty.Name},
		CompletedCourses: completed, IsActive: item.IsActive,
		PersonalDataConsent: item.PersonalDataConsent, ReviewsCount: item.ReviewsCount,
	}
}
