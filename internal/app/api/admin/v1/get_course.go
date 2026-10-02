package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) GetCourse(ctx context.Context, req *adminpb.GetCourseRequest) (*adminpb.GetCourseResponse, error) {
	v, err := s.module.Actions.GetCourse.Execute(ctx, req.GetId())
	return &adminpb.GetCourseResponse{Course: course(v)}, rpcError(err)
}
