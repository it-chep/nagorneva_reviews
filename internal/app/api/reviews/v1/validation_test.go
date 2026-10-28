package reviewsv1

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"testing"
)

func TestValidateLocation(t *testing.T) {
	for _, v := range [][3]float64{{91, 0, 0}, {0, 181, 0}, {0, 0, -1}, {math.NaN(), 0, 0}, {0, math.Inf(1), 0}, {0, 0, math.Inf(1)}} {
		if status.Code(validateLocation(v[0], v[1], v[2])) != codes.InvalidArgument {
			t.Fatalf("accepted %v", v)
		}
	}
	for _, v := range [][3]float64{{0, 0, 0}, {90, 180, 10}, {-90, -180, 20000}} {
		if err := validateLocation(v[0], v[1], v[2]); err != nil {
			t.Fatal(err)
		}
	}
}
func TestPublicContextErrors(t *testing.T) {
	for _, v := range []struct {
		err  error
		code codes.Code
	}{{context.Canceled, codes.Canceled}, {context.DeadlineExceeded, codes.DeadlineExceeded}, {fmt.Errorf("private database details"), codes.Internal}} {
		err := publicError(context.Background(), fmt.Errorf("wrapped: %w", v.err))
		if status.Code(err) != v.code {
			t.Fatal(err)
		}
	}
}
