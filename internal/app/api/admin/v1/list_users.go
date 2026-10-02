package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) ListUsers(ctx context.Context, _ *adminpb.ListUsersRequest) (*adminpb.ListUsersResponse, error) {
	items, err := s.module.Actions.ListUsers.Execute(ctx)
	return &adminpb.ListUsersResponse{Users: users(items)}, rpcError(err)
}
