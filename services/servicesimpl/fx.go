package servicesimpl

import "go.uber.org/fx"

var Module = fx.Module("servicesimpl",
	fx.Provide(NewChatSessionService),
	fx.Provide(NewChatSession),
)
