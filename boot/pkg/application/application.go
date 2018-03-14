package application

import (
	"github.com/kataras/iris"
	log "github.com/kataras/golog"
	"github.com/hi-devops-io/hi/boot/pkg/config"
	"fmt"
)

type Boot struct {
	app *iris.Application
	config *config.AppConfig
}

var (
	boot Boot
)

func init() {
	boot.config = config.BuildAppConfig()
	log.SetLevel(boot.config.Logging.Level)
	log.Debug("application init")
	log.Debug(boot.config)

	boot.app = iris.New()
}

func Instance() *iris.Application {
	return boot.app
}

func Config() *config.AppConfig  {
	return boot.config
}

func Run() {
	serverPort := fmt.Sprintf(":%v", boot.config.Server.Port)
	boot.app.Run(iris.Addr(fmt.Sprintf(serverPort)), iris.WithCharset("UTF-8"), iris.WithoutVersionChecker)
}
