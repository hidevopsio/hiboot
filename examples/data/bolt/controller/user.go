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

package controller

import (
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/examples/data/bolt/entity"
	"github.com/hidevopsio/hiboot/examples/data/bolt/service"
)

//hi: RestController
type UserController struct {
	web.Controller
	userService *service.UserService
}

func init() {
	web.RestController(new(UserController))
}

// Init inject userService automatically
func (c *UserController) Init(userService *service.UserService) {
	c.userService = userService
}

// Post /user
func (c *UserController) Post(user *entity.User) (model.Response, error) {
	err := c.userService.AddUser(user)
	response := new(model.BaseResponse)
	response.SetData(user)
	return response, err
}

// Get /user/{id}
func (c *UserController) GetById(id string) (model.Response, error) {
	user, err := c.userService.GetUser(id)
	response := new(model.BaseResponse)
	if err != nil {
		response.SetCode(http.StatusNotFound)
	} else {
		response.SetData(user)
	}
	return response, err
}

// Delete /user/{id}
func (c *UserController) DeleteById(id string) (response model.Response, err error) {
	err = c.userService.DeleteUser(id)
	response = new(model.BaseResponse)
	return
}
