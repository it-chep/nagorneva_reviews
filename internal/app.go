package internal

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"

	appserver "github.com/nagorneva/nagorneva_reviews/internal/app/server"
	adminmodule "github.com/nagorneva/nagorneva_reviews/internal/module/admin"
	reviewsmodule "github.com/nagorneva/nagorneva_reviews/internal/module/reviews"
)

// modules aggregates application modules in the same way as medblogers_base.
type modules struct {
	admin   *adminmodule.Module
	reviews *reviewsmodule.Module
}

type App struct {
	server   *appserver.Server
	shutdown func()
	modules  modules
}

func New(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initConfig(ctx).initPostgres(ctx).initModules(ctx).initTransport(ctx).err; err != nil {
		return nil, err
	}
	return a, nil
}

// Run owns the application lifecycle; Server.Run starts the HTTP gateway and
// native gRPC listener and gracefully stops both when ctx is cancelled.
func (a *App) Run(ctx context.Context) (runErr error) {
	if a == nil || a.server == nil {
		return fmt.Errorf("application servers are not initialized")
	}
	if a.shutdown != nil {
		defer a.shutdown()
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			runErr = fmt.Errorf("application panic: %v\n%s", recovered, debug.Stack())
		}
	}()
	slog.Info("nagorneva_reviews started", "http", a.server.HTTPAddr(), "grpc", a.server.GRPCAddr())
	return a.server.Run(ctx)
}
