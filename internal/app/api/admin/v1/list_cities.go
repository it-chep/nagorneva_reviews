package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) ListCities(ctx context.Context, _ *adminpb.ListCitiesRequest) (*adminpb.ListCitiesResponse, error) {
	items, err := s.module.Actions.ListCities.Execute(ctx)
	return &adminpb.ListCitiesResponse{Cities: cities(items)}, rpcError(err)
}
