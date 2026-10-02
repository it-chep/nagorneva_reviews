package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) DeleteUser(ctx context.Context, req *adminpb.DeleteUserRequest) (*adminpb.DeleteUserResponse, error) {
	return &adminpb.DeleteUserResponse{}, rpcError(s.module.Actions.DeleteUser.Execute(ctx, req.GetId()))
}
