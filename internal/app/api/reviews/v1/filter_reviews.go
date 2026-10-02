package reviewsv1

import (
	"context"

	reviewspb "github.com/nagorneva/nagorneva_reviews/gen/go/reviews"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) FilterReviews(ctx context.Context, req *reviewspb.FilterReviewsRequest) (*reviewspb.FilterReviewsResponse, error) {
	if (req.Lat == nil) != (req.Lon == nil) {
		return nil, status.Error(codes.InvalidArgument, "lat and lon must be provided together")
	}
	reviews, err := s.module.Actions.FilterReviews.Execute(ctx, req.CourseId, req.Lat, req.Lon, req.GetRadiusKm())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}
	out := make([]*reviewspb.Review, 0, len(reviews))
	for _, review := range reviews {
		out = append(out, &reviewspb.Review{Id: review.ID, DoctorId: review.DoctorID, Rating: int32(review.Rating), Comment: review.Comment, CourseId: review.CourseID, CourseName: review.CourseName, CreatedAt: timestamppb.New(review.CreatedAt)})
	}
	return &reviewspb.FilterReviewsResponse{Reviews: out}, nil
}
