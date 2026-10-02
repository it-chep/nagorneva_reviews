package reviewsv1

import (
	"context"

	reviewspb "github.com/nagorneva/nagorneva_reviews/gen/go/reviews"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) GetDoctorsOnMap(ctx context.Context, req *reviewspb.GetDoctorsOnMapRequest) (*reviewspb.GetDoctorsOnMapResponse, error) {
	doctors, err := s.module.Actions.GetDoctorsOnMap.Execute(ctx, req.GetLat(), req.GetLon(), req.GetRadiusKm())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}
	out := make([]*reviewspb.DoctorOnMap, 0, len(doctors))
	for _, doctor := range doctors {
		item := &reviewspb.DoctorOnMap{Id: doctor.ID, Name: doctor.Name, Photo: doctor.Photo, City: doctor.City, Lat: doctor.Lat, Lon: doctor.Lon, Reviews: make([]*reviewspb.Review, 0, len(doctor.Reviews))}
		for _, review := range doctor.Reviews {
			item.Reviews = append(item.Reviews, &reviewspb.Review{Id: review.ID, DoctorId: review.DoctorID, Rating: int32(review.Rating), Comment: review.Comment, CourseId: review.CourseID, CourseName: review.CourseName, CreatedAt: timestamppb.New(review.CreatedAt)})
		}
		out = append(out, item)
	}
	return &reviewspb.GetDoctorsOnMapResponse{Doctors: out}, nil
}
