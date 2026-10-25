// Package admin composes administrative use cases and their dependencies.
package admin

import (
	"context"

	"github.com/nagorneva/nagorneva_reviews/internal/config"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/action/login"
	s3gateway "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/client/s3"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/action/upload_doctor_photo"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Module struct{ Actions *entities.Aggregator }

func New(ctx context.Context, db postgres.DB, token login.TokenIssuer, cfg config.Config) (*Module, error) {
	client, err := s3gateway.NewGateway(ctx, cfg)
	if err != nil {
		return nil, err
	}
	var images upload_doctor_photo.ImageStorage
	if client != nil {
		images = client
	}
	actions := entities.New(db, token, images)
	if err := actions.Users.EnsureInitialUser.Execute(ctx, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		return nil, err
	}
	return &Module{Actions: actions}, nil
}
