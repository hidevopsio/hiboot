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
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"strings"
)


type Bar struct {
	Greeting string
}


type BarController struct{
	web.JwtController
}

func init()  {
	web.Add(new(BarController))
}

func (c *BarController) GetSayHello(ctx *web.Context)  {

	// decrypt jwt token
	ti := ctx.Values().Get("jwt")
	if ti == nil {
		ctx.ResponseError("failed", http.StatusInternalServerError)
		return
	}
	token := ti.(*jwt.Token)

	var username, password string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		username = c.ParseToken(claims, "username")
		password = c.ParseToken(claims, "password")

		log.Debugf("username: %v, password: %v", username, strings.Repeat("*", len(password)))
	}

	log.Debug("BarController.SayHello")
	language := ctx.Values().GetString(ctx.Application().ConfigurationReadOnly().GetTranslateLanguageContextKey())
	log.Debug(language)
	ctx.ResponseBody("success", &Bar{Greeting: "hello bar"})

}
