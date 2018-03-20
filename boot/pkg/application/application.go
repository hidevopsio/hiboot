package application

import (
	"github.com/kataras/iris"
	log "github.com/kataras/golog"
	"github.com/hidevopsio/hi/boot/pkg/system"
	"fmt"
	"github.com/kataras/iris/context"
)

type Boot struct {
	app *iris.Application
	config *system.Configuration
}

type Health struct {
	Status string `json:"status"`
}

var (
	boot Boot
)

func init() {
	boot.config = system.Build()
	log.SetLevel(boot.config.Logging.Level)
	log.Debug("application init")
	log.Debug(boot.config)

	boot.app = iris.New()

	boot.app.Get("/health", func(ctx context.Context) {
		health := Health{
			Status: "UP",
		}
		ctx.JSON(health)
	})
}

func Instance() *iris.Application {
	return boot.app
}

func Config() *system.Configuration {
	return boot.config
}

func Run() {
	serverPort := fmt.Sprintf(":%v", boot.config.Server.Port)
	boot.app.Run(iris.Addr(fmt.Sprintf(serverPort)), iris.WithCharset("UTF-8"), iris.WithoutVersionChecker)
}
