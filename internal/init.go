package internal

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	adminapi "github.com/nagorneva/nagorneva_reviews/internal/app/api/admin/v1"
	reviewsapi "github.com/nagorneva/nagorneva_reviews/internal/app/api/reviews/v1"
	appserver "github.com/nagorneva/nagorneva_reviews/internal/app/server"
	"github.com/nagorneva/nagorneva_reviews/internal/config"
	adminmodule "github.com/nagorneva/nagorneva_reviews/internal/module/admin"
	reviewsmodule "github.com/nagorneva/nagorneva_reviews/internal/module/reviews"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/auth"
)

type initializer struct {
	*App
	err   error
	cfg   config.Config
	db    *pgxpool.Pool
	token auth.Service
}

func (a *App) initConfig(_ context.Context) *initializer {
	i := &initializer{App: a}
	i.cfg, i.err = config.Load()
	if i.err == nil {
		i.token = auth.New(i.cfg.JWTSecret, i.cfg.JWTTTL)
	}
	return i
}

func (i *initializer) initPostgres(ctx context.Context) *initializer {
	if i.err != nil {
		return i
	}
	i.db, i.err = pgxpool.New(ctx, i.cfg.DatabaseURL)
	if i.err == nil {
		i.shutdown = i.db.Close
		i.err = i.db.Ping(ctx)
	}
	return i
}

func (i *initializer) initModules(ctx context.Context) *initializer {
	if i.err != nil {
		return i
	}
	i.modules.admin, i.err = adminmodule.New(ctx, i.db, i.token, i.cfg)
	if i.err != nil {
		return i
	}
	i.modules.reviews = reviewsmodule.New(i.db)
	return i
}

// initTransport only composes API services. Router and gRPC setup belong to
// app/server, leaving this initializer free of transport implementation.
func (i *initializer) initTransport(ctx context.Context) *initializer {
	if i.err != nil {
		return i
	}
	adminService := adminapi.New(i.modules.admin)
	reviewsService := reviewsapi.New(i.modules.reviews)
	i.server, i.err = appserver.New(ctx, appserver.Config{
		HTTPAddr: i.cfg.HTTPAddr,
		GRPCAddr: i.cfg.GRPCAddr,
		Token:    i.token,
		Debug:    i.cfg.Debug,
	}, adminService, reviewsService)
	return i
}
