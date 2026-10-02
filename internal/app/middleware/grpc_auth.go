package middleware

import (
	"context"
	"strings"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCAdminAuth protects all AdminService RPCs except Login. Authorization is
// deliberately bypassed when debug is enabled.
func GRPCAdminAuth(tokens TokenVerifier, debug bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if debug {
			return handler(ctx, req)
		}
		if !strings.HasPrefix(info.FullMethod, "/nagorneva_reviews.admin.v1.AdminService/") || info.FullMethod == adminpb.AdminService_Login_FullMethodName {
			return handler(ctx, req)
		}
		values := metadata.ValueFromIncomingContext(ctx, "authorization")
		if len(values) != 1 || !strings.HasPrefix(values[0], "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "authorization required")
		}
		if _, err := tokens.Verify(strings.TrimPrefix(values[0], "Bearer ")); err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		return handler(ctx, req)
	}
}
