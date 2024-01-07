package proxy

import (
	"chatrabbit/pkg/services/proxyserv"
	"context"
	"fmt"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/kataras/iris/v12/mvc"
	. "github.com/smartystreets/goconvey/convey"
)

func newApp() *iris.Application {
	app := iris.New()
	mvc.Configure(app.Party("/proxy"), proxyApp)
	return app
}

func proxyApp(app *mvc.Application) {
	service := proxyserv.NewProxyService()
	ctx := context.Background()
	app.Register(ctx, service)
	app.Handle(new(ProxyController))
}

func TestProxy(t *testing.T) {
	app := newApp()
	h := httptest.New(t, app)

	token := "Bearer " + ""
	queryStr := "/v1/completion?query=hello&region=CN&lang=zh"

	Convey("Test Proxy", t, func() {
		h1 := h.GET("/proxy/").
			WithHeader("Authorization", token).
			WithQueryString(queryStr).
			Expect().Status(httptest.StatusOK)
		expectedResp := "\"code\":0,"
		h1.Body().Contains(expectedResp)
		bs := h1.Body().Raw()
		fmt.Println(bs)
	})
}
