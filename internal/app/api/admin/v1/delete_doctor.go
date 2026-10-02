package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) DeleteDoctor(ctx context.Context, req *adminpb.DeleteDoctorRequest) (*adminpb.DeleteDoctorResponse, error) {
	return &adminpb.DeleteDoctorResponse{}, rpcError(s.module.Actions.DeleteDoctor.Execute(ctx, req.GetId()))
}
