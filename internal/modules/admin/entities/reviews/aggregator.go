package reviews

import (
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/action/create_review"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/action/delete_review"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/action/get_doctor_reviews"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/action/get_review"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/action/list_reviews"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/action/update_review"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/dal"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Aggregator struct {
	CreateReview     *create_review.Action
	GetReview        *get_review.Action
	ListReviews      *list_reviews.Action
	UpdateReview     *update_review.Action
	DeleteReview     *delete_review.Action
	GetDoctorReviews *get_doctor_reviews.Action
}

func New(db postgres.DB) *Aggregator {
	repository := dal.New(db)
	return &Aggregator{
		CreateReview:     create_review.New(repository),
		GetReview:        get_review.New(repository),
		ListReviews:      list_reviews.New(repository),
		UpdateReview:     update_review.New(repository),
		DeleteReview:     delete_review.New(repository),
		GetDoctorReviews: get_doctor_reviews.New(repository),
	}
}
