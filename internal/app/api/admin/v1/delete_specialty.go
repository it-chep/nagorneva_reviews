package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) DeleteSpecialty(ctx context.Context, req *adminpb.DeleteSpecialtyRequest) (*adminpb.DeleteSpecialtyResponse, error) {
	return &adminpb.DeleteSpecialtyResponse{}, rpcError(s.module.Actions.Specialties.DeleteSpecialty.Execute(ctx, req.GetId()))
}
