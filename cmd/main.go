package main

import (
	"FantasticLife/api"
	"FantasticLife/config"
	"FantasticLife/server/serverimpl"
	"FantasticLife/services/servicesimpl"
	"FantasticLife/utils"
	"go.uber.org/fx"
)

// TODO: 2. 加入对话接口
// TODO: 3. 缝合im系统
// TODO: 4. 部署验证
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
