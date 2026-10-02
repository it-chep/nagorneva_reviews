package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
	updatedoctor "github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/update_doctor"
)

func (s *Service) UpdateDoctor(ctx context.Context, req *adminpb.UpdateDoctorRequest) (*adminpb.UpdateDoctorResponse, error) {
	v, err := s.module.Actions.UpdateDoctor.Execute(ctx, req.GetId(), updatedoctor.Input{
		Name: req.Name, FullName: req.FullName, CityID: req.CityId, SpecialtyID: req.SpecialtyId, Lat: req.Lat, Lon: req.Lon,
		PersonalDataConsent: req.PersonalDataConsent, IsActive: req.IsActive,
	})
	return &adminpb.UpdateDoctorResponse{Doctor: doctor(v)}, rpcError(err)
}
