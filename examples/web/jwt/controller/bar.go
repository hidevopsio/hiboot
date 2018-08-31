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
	"strings"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/jwt"
	"github.com/hidevopsio/hiboot/pkg/model"
	"net/http"
)


type Bar struct {
	Greeting string
}


type barController struct{
	jwt.Controller
}

func init()  {
	web.RestController(new(barController))
}


func (c *barController) Get() (response model.Response) {
	username := c.JwtProperty("username")
	password := c.JwtProperty("password")
	log.Debugf("username: %v, password: %v", username, strings.Repeat("*", len(password)))

	log.Debug("BarController.SayHello")

	response = new(model.BaseResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("success")
	response.SetData(&Bar{Greeting: "Hello " + username})
	return
}
