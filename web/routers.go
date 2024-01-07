package web

import (
	"chatrabbit/pkg/controllers/proxy"
	"chatrabbit/pkg/services/proxyserv"
	"context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// ProxyRouter 注册路由
func SetupRoutes(app *iris.Application) {

	// API
	proxyApi := app.Party("/proxy")
	{
		mvc.Configure(proxyApi.Party("/{any:path}"), proxyApp)
	}
}

func proxyApp(app *mvc.Application) {

	service := proxyserv.NewProxyService()
	ctx := context.Background()
	app.Register(ctx, service)
	app.Handle(new(proxy.ProxyController))
}
