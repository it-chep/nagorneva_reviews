package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) UpdateUser(ctx context.Context, req *adminpb.UpdateUserRequest) (*adminpb.UpdateUserResponse, error) {
	v, err := s.module.Actions.UpdateUser.Execute(ctx, req.GetId(), req.Email, req.Password)
	return &adminpb.UpdateUserResponse{User: user(v)}, rpcError(err)
}
