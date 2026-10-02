package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) GetSpecialty(ctx context.Context, req *adminpb.GetSpecialtyRequest) (*adminpb.GetSpecialtyResponse, error) {
	v, err := s.module.Actions.GetSpecialty.Execute(ctx, req.GetId())
	return &adminpb.GetSpecialtyResponse{Specialty: specialty(v)}, rpcError(err)
}
