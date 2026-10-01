package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	server   *http.Server
	shutdown func()
}

func New(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initConfig(ctx).initPostgres(ctx).initS3(ctx).initModules(ctx).initRouter(ctx).err; err != nil {
		return nil, err
	}
	return a, nil
}
func Run(ctx context.Context) error {
	a, err := New(ctx)
	if err != nil {
		return err
	}
	defer a.shutdown()
	errCh := make(chan error, 1)
	go func() {
		slog.Info("nagorneva_reviews started", "http", a.server.Addr)
		errCh <- a.server.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return a.server.Shutdown(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("http server: %w", err)
	}
}
