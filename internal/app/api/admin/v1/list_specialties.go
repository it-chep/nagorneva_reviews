package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) ListSpecialties(ctx context.Context, _ *adminpb.ListSpecialtiesRequest) (*adminpb.ListSpecialtiesResponse, error) {
	items, err := s.module.Actions.ListSpecialties.Execute(ctx)
	return &adminpb.ListSpecialtiesResponse{Specialties: specialties(items)}, rpcError(err)
}
