package main

import (
	"chatrabbit/config"
	"chatrabbit/config/common"
	"chatrabbit/pkg/infra/log"
	"chatrabbit/web"
	"flag"
	"fmt"
	"os"

	"github.com/kataras/iris/v12"
)

func newApp() *iris.Application {
	app := iris.New()

	// iris log level
	logLevel := config.GetString(common.LOG_LEVEL)
	app.Logger().SetLevel(logLevel)

	web.SetupRoutes(app)

	// Iris log
	f, _ := os.Create("iris.log")
	app.Logger().SetOutput(f)

	return app
}

func InitModule() {
	// database.Init()
}

func main() {
	log.ConfigLogger()
	log.Info("start proxy server")

	configFile := flag.String("config", "", "specify config file name")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}
	ret := config.InitConfig(*configFile)
	if !ret {
		flag.Usage()
		os.Exit(-1)
	}
	InitModule()

	logLevel := config.GetString(common.LOG_LEVEL)
	// project log level
	log.SetLogLevel(logLevel)

	app := newApp()

	// iris config
	c := iris.Configuration{
		DisableStartupLog:                 false,
		DisableBodyConsumptionOnUnmarshal: true,
		Charset:                           "UTF-8",
	}
	irisConfig := iris.WithConfiguration(c)
	addr := fmt.Sprintf("%s:%s", config.GetString(common.SERVER_HOST), config.GetString(common.SERVER_PORT))
	log.Infof("server run at %v", addr)
	app.Run(iris.Addr(addr), irisConfig)
}
