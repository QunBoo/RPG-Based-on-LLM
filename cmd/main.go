package main

import (
	"FantasticLife/api"
	"FantasticLife/config"
	"FantasticLife/server/serverimpl"
	"FantasticLife/services/servicesimpl"
	"FantasticLife/utils"

	"go.uber.org/fx"
)

// TODO：实现Websocket协议成功升级，从http升级为websocket
// TODO：实现websocket的ping-pong机制，保持连接
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
