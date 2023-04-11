package routers

import (
	"chatrabbit/pkg/controllers/cproxy"
	"chatrabbit/pkg/services/sproxy"
	"context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// ProxyRouter 注册路由
func RegistProxyRouter(app *iris.Application) {

	// API
	proxyApi := app.Party("/proxy")
	{
		mvc.Configure(proxyApi.Party("/{any:path}"), proxyApp)
	}
}

func proxyApp(app *mvc.Application) {

	service := sproxy.NewProxyService()
	ctx := context.Background()
	app.Register(ctx, service)
	app.Handle(new(cproxy.ProxyController))
}
