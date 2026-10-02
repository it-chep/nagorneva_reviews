package adminv1

import (
	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
	adminmodule "github.com/nagorneva/nagorneva_reviews/internal/module/admin"
)

// Service adapts generated gRPC methods to admin actions.
type Service struct {
	adminpb.UnimplementedAdminServiceServer
	module *adminmodule.Module
}

// New creates the admin gRPC service.
func New(module *adminmodule.Module) *Service {
	return &Service{module: module}
}
