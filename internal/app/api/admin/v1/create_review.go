package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) CreateReview(ctx context.Context, req *adminpb.CreateReviewRequest) (*adminpb.CreateReviewResponse, error) {
	v, err := s.module.Actions.CreateReview.Execute(ctx, req.GetDoctorId(), req.GetRating(), req.GetComment(), req.GetCourseId(), req.GetIsActive())
	return &adminpb.CreateReviewResponse{Review: review(v)}, rpcError(err)
}
