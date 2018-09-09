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

package logging

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/middleware/logger"
)

type configuration struct {
	app.PreConfiguration
	Properties Properties `mapstructure:"logging"`
	applicationContext app.ApplicationContext
}

func newConfiguration(applicationContext app.ApplicationContext) *configuration {
	return &configuration{
		applicationContext: applicationContext,
	}
}

func init() {
	app.AutoConfiguration(newConfiguration)
}

// LoggerHandler config logger handler
func (c *configuration) LoggerHandler() context.Handler {
	loggerHandler := logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
		// Query appends the url query to the Path.
		//Query: true,

		//Columns: true,

		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
		// will be added to the logs.
		MessageContextKeys: []string{"logger_message"},

		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
		MessageHeaderKeys: []string{"User-Agent"},
	})

	c.applicationContext.Use(loggerHandler)

	return loggerHandler
}