package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) UpdateSpecialty(ctx context.Context, req *adminpb.UpdateSpecialtyRequest) (*adminpb.UpdateSpecialtyResponse, error) {
	v, err := s.module.Actions.UpdateSpecialty.Execute(ctx, req.GetId(), req.Name)
	return &adminpb.UpdateSpecialtyResponse{Specialty: specialty(v)}, rpcError(err)
}
