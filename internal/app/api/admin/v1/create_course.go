package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) CreateCourse(ctx context.Context, req *adminpb.CreateCourseRequest) (*adminpb.CreateCourseResponse, error) {
	v, err := s.module.Actions.CreateCourse.Execute(ctx, req.GetName(), req.GetSiteLink())
	return &adminpb.CreateCourseResponse{Course: course(v)}, rpcError(err)
}
