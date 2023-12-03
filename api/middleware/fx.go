package middleware

import "go.uber.org/fx"

var Module = fx.Module("middle",
	fx.Provide(
		fx.Annotate(CORS, fx.ResultTags(`name:"cors"`)),
		fx.Annotate(ZapLogger, fx.ResultTags(`name:"zaplogger"`)),
	),
)
