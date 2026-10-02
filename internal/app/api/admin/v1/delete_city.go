package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) DeleteCity(ctx context.Context, req *adminpb.DeleteCityRequest) (*adminpb.DeleteCityResponse, error) {
	return &adminpb.DeleteCityResponse{}, rpcError(s.module.Actions.DeleteCity.Execute(ctx, req.GetId()))
}
