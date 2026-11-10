package adminv1

import (
	"context"
	"testing"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFilterDoctorsRejectsInvalidFilters(t *testing.T) {
	for _, req := range []*adminpb.FilterDoctorsRequest{
		{CityIds: []int64{0}}, {SpecialtyIds: []int64{-1}}, {CourseIds: []int64{1, 0}},
		{ReviewsSort: adminpb.ReviewsSort(-1)}, {ReviewsSort: adminpb.ReviewsSort(3)},
	} {
		_, err := (&Service{}).FilterDoctors(context.Background(), req)
		if status.Code(err) != codes.InvalidArgument {
			t.Fatalf("request %v: %v", req, err)
		}
	}
}
