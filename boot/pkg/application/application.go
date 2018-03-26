// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


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
