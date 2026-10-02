package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
	createdoctor "github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/create_doctor"
)

func (s *Service) CreateDoctor(ctx context.Context, req *adminpb.CreateDoctorRequest) (*adminpb.CreateDoctorResponse, error) {
	v, err := s.module.Actions.CreateDoctor.Execute(ctx, createdoctor.Input{
		Name: req.GetName(), FullName: req.GetFullName(), CityID: req.GetCityId(), SpecialtyID: req.GetSpecialtyId(),
		Lat: req.GetLat(), Lon: req.GetLon(), PersonalDataConsent: req.GetPersonalDataConsent(), IsActive: req.GetIsActive(),
	})
	return &adminpb.CreateDoctorResponse{Doctor: doctor(v)}, rpcError(err)
}
