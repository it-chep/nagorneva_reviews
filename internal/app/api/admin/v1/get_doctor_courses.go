package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) GetDoctorCourses(ctx context.Context, req *adminpb.GetDoctorCoursesRequest) (*adminpb.GetDoctorCoursesResponse, error) {
	items, err := s.module.Actions.Doctors.GetDoctorCourses.Execute(ctx, req.GetId())
	if err != nil {
		return nil, rpcError(err)
	}
	return &adminpb.GetDoctorCoursesResponse{Courses: courses(items)}, nil
}
