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

// Package websocket provides web socket auto configuration for web/cli application
package websocket

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/kataras/iris/websocket"
)

// Profile websocket profile name
const Profile = "websocket"

type configuration struct {
	at.AutoConfiguration

	Properties properties `mapstructure:"websocket"`
}

func newConfiguration() *configuration {
	return &configuration{}
}

func init() {
	app.Register(newConfiguration)
}

// Server websocket server
func (c *configuration) Server() *websocket.Server {
	ws := websocket.New(websocket.Config{
		ReadBufferSize:  c.Properties.ReadBufferSize,
		WriteBufferSize: c.Properties.WriteBufferSize,
	})
	return ws
}

// Upgrade websocket connection
func (c *configuration) Connection(ctx context.Context, server *websocket.Server) websocket.Connection {
	return server.Upgrade(ctx)
}
