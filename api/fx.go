package api

import (
	"FantasticLife/api/middleware"
	"FantasticLife/config"
	"FantasticLife/server/serverimpl/WebSocket"
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Servers struct {
	fx.In

	Server1 *http.Server `name:"server1"`
	// 如果需要更多的服务器，可以继续在这里添加
}

func HttpServerLifecycle(lc fx.Lifecycle, servers Servers, logger *zap.Logger, config *config.Config, clientManager *WebSocket.ClientManager) {
	server1 := servers.Server1
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("http server starting", zap.String("address", server1.Addr))
				if err := server1.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("server failed to start", zap.Error(err))
				}
			}()
			go WebSocket.StartWebSocket(config, clientManager, logger)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			if err := server1.Shutdown(_ctx); err != nil {
				logger.Error("server shutdown failed", zap.Error(err))
			} else {
				logger.Info("server gracefully shutdown")
			}
			return nil
		},
	})
}

// TODO：写第二个server，启动websocket服务器程序，设定webSocketPort
var Module = fx.Module("router",
	middleware.Module,
	fx.Provide(
		fx.Annotated{
			Name:   "server1",
			Target: NewHttpServer,
		},
	),
	fx.Invoke(HttpServerLifecycle),
)
