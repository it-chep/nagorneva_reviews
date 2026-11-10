package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) DeleteDoctorCourse(ctx context.Context, req *adminpb.DeleteDoctorCourseRequest) (*adminpb.DeleteDoctorCourseResponse, error) {
	if err := s.module.Actions.Doctors.DeleteDoctorCourse.Execute(ctx, req.GetId(), req.GetCourseId()); err != nil {
		return nil, rpcError(err)
	}
	return &adminpb.DeleteDoctorCourseResponse{}, nil
}
