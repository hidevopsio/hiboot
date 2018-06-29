package controllers

import (
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/examples/db/bolt/services"
	"github.com/hidevopsio/hiboot/examples/db/bolt/domain"
)

type UserController struct {
	web.Controller

	UserService *services.UserService `component:"service"`
}

func init() {
	web.Add(new(UserController))
}

func (c *UserController) Post(ctx *web.Context) {

	user := &domain.User{}
	err := ctx.RequestBody(user)
	if err == nil {
		c.UserService.AddUser(user)

		ctx.ResponseBody("success", user)
	}

}

func (c *UserController) Get(ctx *web.Context) {

	id := ctx.URLParam("id")

	user, err := c.UserService.GetUser(id)
	if err != nil {
		ctx.ResponseError("Resource is not found", http.StatusNotFound)
	} else {
		ctx.ResponseBody("success", user)
	}
}

func (c *UserController) Delete(ctx *web.Context) {

	id := ctx.URLParam("id")

	err := c.UserService.DeleteUser(id)
	if err != nil {
		ctx.ResponseError("Failed", http.StatusInternalServerError)
	} else {
		ctx.ResponseBody("success", nil)
	}
}
