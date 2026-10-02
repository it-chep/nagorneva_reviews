package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) GetDoctor(ctx context.Context, req *adminpb.GetDoctorRequest) (*adminpb.GetDoctorResponse, error) {
	v, err := s.module.Actions.GetDoctor.Execute(ctx, req.GetId())
	return &adminpb.GetDoctorResponse{Doctor: doctor(v)}, rpcError(err)
}
