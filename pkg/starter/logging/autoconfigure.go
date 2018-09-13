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
	Properties         Properties `mapstructure:"logging"`
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
		Status: c.Properties.Status,
		// IP displays request's remote address
		IP: c.Properties.IP,
		// Method displays the http method
		Method: c.Properties.Method,
		// Path displays the request path
		Path: c.Properties.Path,
		// Query appends the url query to the Path.
		Query: c.Properties.Query,

		Columns: c.Properties.Columns,

		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
		// will be added to the logs.
		MessageContextKeys: c.Properties.ContextKeys,

		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
		MessageHeaderKeys: c.Properties.HeaderKeys,
	})

	c.applicationContext.Use(loggerHandler)

	return loggerHandler
}
