package controller

import (
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/examples/data/bolt/model"
	"github.com/hidevopsio/hiboot/examples/data/bolt/service"
)

//hi: RestController
type UserController struct {
	web.Controller

	UserService *service.UserService `inject:""`
}

func init() {
	web.Add(new(UserController))
}

//hi: Constructor Injection
func (c *UserController) Init(userService *service.UserService) {
	c.UserService = userService
}

//hi: method=POST
func (c *UserController) Post() {

	user := &model.User{}
	err := c.Ctx.RequestBody(user)
	if err == nil {
		c.UserService.AddUser(user)

		c.Ctx.ResponseBody("success", user)
	}
}

func (c *UserController) Get() {

	id := c.Ctx.URLParam("id")

	user, err := c.UserService.GetUser(id)
	if err != nil {
		c.Ctx.ResponseError(err.Error(), http.StatusNotFound)
	} else {
		c.Ctx.ResponseBody("success", user)
	}
}

func (c *UserController) Delete() {

	id := c.Ctx.URLParam("id")

	err := c.UserService.DeleteUser(id)
	if err != nil {
		c.Ctx.ResponseError(err.Error(), http.StatusInternalServerError)
	} else {
		c.Ctx.ResponseBody("success", nil)
	}
}
