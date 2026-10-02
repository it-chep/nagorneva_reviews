package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) CreateSpecialty(ctx context.Context, req *adminpb.CreateSpecialtyRequest) (*adminpb.CreateSpecialtyResponse, error) {
	v, err := s.module.Actions.CreateSpecialty.Execute(ctx, req.GetName())
	return &adminpb.CreateSpecialtyResponse{Specialty: specialty(v)}, rpcError(err)
}
