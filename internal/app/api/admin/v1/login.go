package adminv1

import (
	"context"
	"errors"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/login"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Login(ctx context.Context, req *adminpb.LoginRequest) (*adminpb.LoginResponse, error) {
	token, err := s.module.Actions.Login.Execute(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, login.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid email or password")
		}
		return nil, rpcError(err)
	}
	return &adminpb.LoginResponse{AccessToken: token, TokenType: "Bearer"}, nil
}
