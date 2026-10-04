package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) UpdateCourse(ctx context.Context, req *adminpb.UpdateCourseRequest) (*adminpb.UpdateCourseResponse, error) {
	v, err := s.module.Actions.UpdateCourse.Execute(ctx, req.GetId(), req.Name, req.SiteLink)
	return &adminpb.UpdateCourseResponse{Course: course(v)}, rpcError(err)
}
