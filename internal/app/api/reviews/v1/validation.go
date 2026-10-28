package reviewsv1

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
)

func validateLocation(lat, lon, radius float64) error {
	if math.IsNaN(lat) || math.IsInf(lat, 0) || lat < -90 || lat > 90 || math.IsNaN(lon) || math.IsInf(lon, 0) || lon < -180 || lon > 180 || math.IsNaN(radius) || math.IsInf(radius, 0) || radius < 0 {
		return status.Error(codes.InvalidArgument, "invalid coordinates or radius: lat must be [-90,90], lon [-180,180], radius_km finite and non-negative")
	}
	return nil
}
func publicError(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return status.FromContextError(ctx.Err()).Err()
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return status.FromContextError(err).Err()
	}
	return status.Error(codes.Internal, "internal server error")
}
