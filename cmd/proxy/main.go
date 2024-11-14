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

	configServe, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	// iris log level
	logLevel := configServe.GetString(common.LOG_LEVEL)
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
	// configFile := flag.String("config", "", "specify config file name")
	// help := flag.Bool("h", false, "help")
	// flag.Parse()

	// if *help {
	// 	flag.Usage()
	// 	os.Exit(0)
	// }
	_, err := config.GetConfig()
	if err != nil {
		fmt.Printf("init config failed: %v\n", err)
		flag.Usage()
		os.Exit(-1)
	}

	log.ConfigLogger()
	log.Info("start proxy server")

	InitModule()

	configServe, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	logLevel := configServe.GetString(common.LOG_LEVEL)
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
	addr := fmt.Sprintf("%s:%s", configServe.GetString(common.SERVER_HOST), configServe.GetString(common.SERVER_PORT))
	log.Infof("server run at %v", addr)
	app.Run(iris.Addr(addr), irisConfig)
}
