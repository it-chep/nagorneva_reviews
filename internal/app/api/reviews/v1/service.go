package reviewsv1

import (
	reviewspb "github.com/nagorneva/nagorneva_reviews/gen/go/reviews"
	reviewsmodule "github.com/nagorneva/nagorneva_reviews/internal/module/reviews"
)

// Service contains the dependencies shared by public review endpoints.
type Service struct {
	reviewspb.UnimplementedReviewsServiceServer
	module *reviewsmodule.Module
}

func New(module *reviewsmodule.Module) *Service {
	return &Service{module: module}
}
