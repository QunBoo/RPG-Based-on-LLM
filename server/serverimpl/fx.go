package serverimpl

import "go.uber.org/fx"

var Module = fx.Module("serverimpl",
	fx.Provide(NewLLMBOT),
	fx.Provide(NewLLMTransceiver),
)
