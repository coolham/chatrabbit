package web

import (
	"os"

	"github.com/kataras/iris/v12"
)

func InitApp() *iris.Application {
	app := iris.New()

	// app.Use(middleware.GlobalMiddleware) // 全局中间件

	app.Logger().SetLevel("info") // Iris log level

	// Iris log
	f, _ := os.Create("iris.log")
	app.Logger().SetOutput(f)

	// 加载模板文件
	// app.RegisterView(iris.HTML("../../views", ".html"))

	return app
}
