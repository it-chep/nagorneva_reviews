package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) GetUser(ctx context.Context, req *adminpb.GetUserRequest) (*adminpb.GetUserResponse, error) {
	v, err := s.module.Actions.GetUser.Execute(ctx, req.GetId())
	return &adminpb.GetUserResponse{User: user(v)}, rpcError(err)
}
