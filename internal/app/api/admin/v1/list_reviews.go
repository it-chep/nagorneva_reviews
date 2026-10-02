package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) ListReviews(ctx context.Context, _ *adminpb.ListReviewsRequest) (*adminpb.ListReviewsResponse, error) {
	items, err := s.module.Actions.ListReviews.Execute(ctx)
	return &adminpb.ListReviewsResponse{Reviews: reviews(items)}, rpcError(err)
}
