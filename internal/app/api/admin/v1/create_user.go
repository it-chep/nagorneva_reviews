package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) CreateUser(ctx context.Context, req *adminpb.CreateUserRequest) (*adminpb.CreateUserResponse, error) {
	v, err := s.module.Actions.CreateUser.Execute(ctx, req.GetEmail(), req.GetPassword())
	return &adminpb.CreateUserResponse{User: user(v)}, rpcError(err)
}
