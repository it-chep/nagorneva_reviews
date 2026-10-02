package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) CreateCity(ctx context.Context, req *adminpb.CreateCityRequest) (*adminpb.CreateCityResponse, error) {
	v, err := s.module.Actions.CreateCity.Execute(ctx, req.GetName(), req.GetLat(), req.GetLon())
	return &adminpb.CreateCityResponse{City: city(v)}, rpcError(err)
}
