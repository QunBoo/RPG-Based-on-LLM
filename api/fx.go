package api

import (
	"FantasticLife/api/middleware"
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func HttpServerLifecycle(lc fx.Lifecycle, server *http.Server, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("http server starting", zap.String("address", server.Addr))
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("server failed to start", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			if err := server.Shutdown(_ctx); err != nil {
				logger.Error("server shutdown failed", zap.Error(err))
			} else {
				logger.Info("server gracefully shutdown")
			}
			return nil
		},
	})
}

var Module = fx.Module("router",
	middleware.Module,
	fx.Provide(
		NewHttpServer,
	),
	fx.Invoke(HttpServerLifecycle),
)
