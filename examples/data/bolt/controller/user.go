package controller

import (
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
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
	web.Add(new(UserController))
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
