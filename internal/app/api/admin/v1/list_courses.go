package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) ListCourses(ctx context.Context, _ *adminpb.ListCoursesRequest) (*adminpb.ListCoursesResponse, error) {
	items, err := s.module.Actions.ListCourses.Execute(ctx)
	return &adminpb.ListCoursesResponse{Courses: courses(items)}, rpcError(err)
}
