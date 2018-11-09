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
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/websocket"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
)

const (
	// Profile websocket profile name
	Profile = "websocket"
	// All is the string which the Emitter use to send a message to all.
	All = ""
	// Broadcast is the string which the Emitter use to send a message to all except this connection.
	Broadcast = ";to;all;except;me;"
)

type configuration struct {
	at.AutoConfiguration

	Properties properties `mapstructure:"websocket"`
}

// Connection is the websocket connection
type Connection interface {
	websocket.Connection
}

type ConnectionFunc func(ctx context.Context) websocket.Connection

// ConnectionFunc is the websocket connection function
type HandlerFunc func(ctx *web.Context, constructor HandlerConstructor) Connection

// ConnectionFunc is the websocket connection function
type HandlerConstructor func(conn Connection) Handler

type Server struct {
	*websocket.Server
}

func newConfiguration() *configuration {
	return &configuration{}
}

func init() {
	app.Register(newConfiguration)
}

// Server websocket server
func (c *configuration) Server() *Server {
	s := websocket.New(websocket.Config{
		ReadBufferSize:  c.Properties.ReadBufferSize,
		WriteBufferSize: c.Properties.WriteBufferSize,
	})

	return &Server{
		Server: s,
	}
}

// ConnectionFunc is func that websocket connection created from
func (c *configuration) ConnectionFunc(server *Server) ConnectionFunc {
	return server.Upgrade
}

// HandlerFunc websocket connection for runtime dependency injection
func (c *configuration) HandlerFunc(connectionFunc ConnectionFunc) HandlerFunc {
	return func(ctx *web.Context, constructor HandlerConstructor) Connection {
		headerConnection := ctx.GetHeader("Connection")
		if headerConnection == "Upgrade" {
			conn := connectionFunc(ctx)
			HandleConnection(constructor, conn)
			return conn
		} else {
			log.Errorf("Header connection is: %v", headerConnection)
		}
		return nil
	}
}
