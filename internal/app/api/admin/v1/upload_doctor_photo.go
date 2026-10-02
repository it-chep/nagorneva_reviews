package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) UploadDoctorPhoto(ctx context.Context, req *adminpb.UploadDoctorPhotoRequest) (*adminpb.UploadDoctorPhotoResponse, error) {
	v, err := s.module.Actions.UploadDoctorPhoto.Execute(ctx, req.GetId(), req.GetPhoto(), req.GetFilename(), req.GetContentType())
	if err != nil {
		return nil, rpcError(err)
	}
	return &adminpb.UploadDoctorPhotoResponse{Doctor: doctor(v)}, nil
}
