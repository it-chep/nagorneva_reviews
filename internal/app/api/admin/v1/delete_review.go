package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
)

func (s *Service) DeleteReview(ctx context.Context, req *adminpb.DeleteReviewRequest) (*adminpb.DeleteReviewResponse, error) {
	return &adminpb.DeleteReviewResponse{}, rpcError(s.module.Actions.DeleteReview.Execute(ctx, req.GetId()))
}
