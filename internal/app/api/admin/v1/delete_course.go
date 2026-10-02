package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) DeleteCourse(ctx context.Context, req *adminpb.DeleteCourseRequest) (*adminpb.DeleteCourseResponse, error) {
	return &adminpb.DeleteCourseResponse{}, rpcError(s.module.Actions.DeleteCourse.Execute(ctx, req.GetId()))
}
