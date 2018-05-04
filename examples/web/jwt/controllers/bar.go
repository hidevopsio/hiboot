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

package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type UserRequest struct {
	Username string
	Password string
}


type Bar struct {
	Greeting string
}


type BarController struct{
	web.Controller
}


func init()  {
	web.Add(&BarController{
		web.Controller{
			ContextMapping: "/bars",
			AuthType:       web.AuthTypeJwt,
		},
	})
}

func (c *BarController) GetSayHello(ctx *web.Context)  {
	log.Debug("BarController.SayHello")

	ctx.Response("Success", &Bar{Greeting: "hello bar"})

}
