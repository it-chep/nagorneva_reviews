package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) ListDoctors(ctx context.Context, _ *adminpb.ListDoctorsRequest) (*adminpb.ListDoctorsResponse, error) {
	items, err := s.module.Actions.ListDoctors.Execute(ctx)
	return &adminpb.ListDoctorsResponse{Doctors: doctors(items)}, rpcError(err)
}
