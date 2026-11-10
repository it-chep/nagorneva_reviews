package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) AddDoctorCourse(ctx context.Context, req *adminpb.AddDoctorCourseRequest) (*adminpb.AddDoctorCourseResponse, error) {
	if err := s.module.Actions.Doctors.AddDoctorCourse.Execute(ctx, req.GetId(), req.GetCourseId()); err != nil {
		return nil, rpcError(err)
	}
	return &adminpb.AddDoctorCourseResponse{}, nil
}
