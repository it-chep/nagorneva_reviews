package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) UpdateCity(ctx context.Context, req *adminpb.UpdateCityRequest) (*adminpb.UpdateCityResponse, error) {
	v, err := s.module.Actions.UpdateCity.Execute(ctx, req.GetId(), req.Name, req.Lat, req.Lon)
	return &adminpb.UpdateCityResponse{City: city(v)}, rpcError(err)
}
