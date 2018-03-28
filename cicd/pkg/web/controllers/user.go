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
	"github.com/kataras/iris"
	"github.com/hidevopsio/hi/cicd/pkg/auth"
	"github.com/hidevopsio/hi/boot/pkg/application"
)

type UserRequest struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Message string                `json:"message"`
	Token   *application.JwtToken `json:"token"`
}

// Operations about object
type UserController struct {
}

// @Title Login
// @Description login
// @Param	body
// @Success 200 {string}
// @Failure 403 body is empty
// @router / [post]
func (c *UserController) Login(ctx iris.Context) {
	var request UserRequest
	var response *UserResponse

	err := ctx.ReadJSON(&request)
	if err != nil {
		ctx.Values().Set("error", "login failed, read and parse request body failed. " + err.Error())
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	// invoke models
	u := &auth.User{}
	token, message, err := u.Login(request.Url, request.Username, request.Password)
	if err == nil {

		jwtToken, err := application.GenerateJwtToken(application.MapJwt{
			"username": request.Username,
			"password": token,
		}, 24)
		if err == nil {
			response = &UserResponse{
				Message: message,
				Token:   &jwtToken,
			}
		} else {
			ctx.Values().Set("error", "login failed, generating token failed. " + err.Error())
			ctx.StatusCode(iris.StatusInternalServerError)
		}
	} else {
		ctx.Values().Set("error", "login failed, wrong username or password " + err.Error())
		ctx.StatusCode(iris.StatusForbidden)
		return
	}

	// just for debug now
	ctx.JSON(response)
}
