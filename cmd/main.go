package main

import (
	"FantasticLife/api"
	"FantasticLife/config"
	"FantasticLife/server/serverimpl"
	"FantasticLife/services/servicesimpl"
	"FantasticLife/utils"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Supply(fx.Annotate(":8080", fx.ResultTags(`name:"hostPort"`))),
		serverimpl.Module,
		servicesimpl.Module,
		api.Module,
		config.Module,
		utils.Module,
	)
	app.Run()
}
