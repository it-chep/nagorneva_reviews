package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) GetDoctorReviews(ctx context.Context, req *adminpb.GetDoctorReviewsRequest) (*adminpb.GetDoctorReviewsResponse, error) {
	items, err := s.module.Actions.GetDoctorReviews.Execute(ctx, req.GetId())
	return &adminpb.GetDoctorReviewsResponse{Reviews: reviews(items)}, rpcError(err)
}
