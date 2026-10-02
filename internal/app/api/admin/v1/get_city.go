package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) GetCity(ctx context.Context, req *adminpb.GetCityRequest) (*adminpb.GetCityResponse, error) {
	v, err := s.module.Actions.GetCity.Execute(ctx, req.GetId())
	return &adminpb.GetCityResponse{City: city(v)}, rpcError(err)
}
