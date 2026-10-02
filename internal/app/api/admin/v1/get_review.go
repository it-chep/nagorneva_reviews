package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) GetReview(ctx context.Context, req *adminpb.GetReviewRequest) (*adminpb.GetReviewResponse, error) {
	v, err := s.module.Actions.GetReview.Execute(ctx, req.GetId())
	return &adminpb.GetReviewResponse{Review: review(v)}, rpcError(err)
}
