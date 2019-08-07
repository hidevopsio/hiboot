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

// Package router provides the web application example that handle customized router in controller
// main package
package main

// import web starter from hiboot
import (
	"fmt"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	"time"
)

// UserController Rest Controller with path /
// RESTful Controller, derived from at.RestController.
type UserController struct {
	// at.RestController or at.RestController must be embedded here
	at.RestController

	// RequestMapping The request mapping of this controller is '/' by default, if you add value tag with value /user,
	// then Hiboot will inject /user to UserController.RequestMapping
	at.RequestMapping `value:"/user"`
}

// UserVO
type UserVO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// User
type User struct {
	ID int `json:"id"`
	UserVO
	model.BaseData
}

// UserRequest
type UserRequest struct {
	at.RequestBody

	UserVO
}

// UserResponse
type UserResponse struct {
	at.ResponseBody
	model.BaseResponse
}

func newUserController() *UserController {
	return &UserController{}
}

// Create
func (c *UserController) CreateUser(
	request *UserRequest,
	m struct {
	at.PostMapping;
	at.Path `value:"/"`
},
) (response *UserResponse, err error) {

	// response
	response = new(UserResponse)
	user := new(User)
	user.ID = 123456
	user.Username = request.Username
	user.Password = request.Password
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsDeleted = false
	response.SetData(user)
	return
}

// Get
func (c *UserController) GetUserById(
	id int,
	name string,
	m struct {
	at.GetMapping
	at.Path `value:"/{id:int}/and/{name}"`
},
) (response *UserResponse, err error) {

	// response
	response = new(UserResponse)
	user := new(User)
	user.ID = id
	user.Username = name
	user.Password = "******"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsDeleted = false
	response.SetData(user)
	return
}

// Patch
func (c *UserController) Patch(
	id int,
	m struct {
	at.PatchMapping
	at.Path `value:"/"`
},
) (response *UserResponse, err error) {

	// response
	response = new(UserResponse)
	return
}

// Delete
func (c *UserController) DeleteUser(
	id int,
	m struct {
	at.DeleteMapping
	at.Path `value:"/{id:int}"`
},
) (response *UserResponse, err error) {

	// response
	response = new(UserResponse)
	return
}

// List
func (c *UserController) ListUser(
	m struct {
	at.GetMapping
	at.Path `value:"/"`
},
) (response *UserResponse, err error) {

	// response
	response = new(UserResponse)
	user := new(User)
	user.ID = 101010
	user.Username = "Donald"
	user.Password = "Trump"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsDeleted = false

	users := []*User{
		user,
	}

	response.SetData(users)
	return
}

type orgController struct {
	at.RestController

	at.RequestMapping `value:"/organization"`
}

// newOrgController is the constructor for orgController
// you may inject dependency to constructor
func newOrgController() *orgController {
	return new(orgController)
}

// GetOfficialSite
func (c *orgController) GetOfficialSite(
	at struct {
	// at.Method is a annotations to define request mapping for http method GET
	at.Method `value:"GET"`
	// at.Path is a annotations to define request mapping for path /official-site
	at.Path `value:"/official-site"`
}) string {

	return "https://hidevops.io"
}

// GetWithPathParamIdAndName
func (c *orgController) GetWithPathParamIdAndName(
	id int,
	name string,
	at struct {
	// at.GetMapping is an annotation to define request mapping for http method GET,
	at.GetMapping
	// at.Path is an annotation to define path for /:id/and/:name
	// As you can see, Hiboot support path param in /:id/and/:name or /{id:int}/and/{name}
	at.Path `value:"/:id/and/:name"`
},
) string {

	return fmt.Sprintf("https://hidevops.io/%v/%v", id, name)
}

// Hiboot main function
func main() {
	// create new web application and run it
	web.NewApplication(newUserController, newOrgController).
		Run()
}
