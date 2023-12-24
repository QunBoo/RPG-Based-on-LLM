package serverimpl

import (
	"FantasticLife/server/serverimpl/WebSocket"
	"FantasticLife/server/serverimpl/cache"
	"go.uber.org/fx"
)

var Module = fx.Module("serverimpl",
	fx.Provide(NewLLMBOT),
	fx.Provide(NewLLMTransceiver),
	fx.Provide(WebSocket.NewClientManager),
	fx.Provide(cache.NewRedisCli),
)
