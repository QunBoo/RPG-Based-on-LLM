package utils

import "go.uber.org/fx"

var Module = fx.Module("utils",
	fx.Provide(
		NewZapLogger,
	),
)
