package adminv1

import (
	"context"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
	updatereview "github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/update_review"
)

func (s *Service) UpdateReview(ctx context.Context, req *adminpb.UpdateReviewRequest) (*adminpb.UpdateReviewResponse, error) {
	v, err := s.module.Actions.UpdateReview.Execute(ctx, req.GetId(), updatereview.Input{
		Rating: req.Rating, Comment: req.Comment, CourseID: req.CourseId, IsActive: req.IsActive,
	})
	return &adminpb.UpdateReviewResponse{Review: review(v)}, rpcError(err)
}
