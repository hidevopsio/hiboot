package controller

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	"net/http"
)

type UserController struct {
	at.RestController

	at.RequestMapping `value:"/user" `
}

func init() {
	app.Register(newUserController)
}

func newUserController() *UserController {
	return &UserController{}
}

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	at.ResponseBody `json:"-"`
	model.BaseResponse
	Data *User `json:"data"`
}

// GetUser
func (c *UserController) GetUser(at struct{
	at.GetMapping `value:"/{id}"`
	at.Operation   `id:"Update Employee" description:"Get User"`
	at.RequiresPermissions `value:"printer:print:hp890"`
}, id int, ctx context.Context) (response *UserResponse) {

	response = new(UserResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("Success")
	response.Data = &User{ID: id, Username: "john.deng", Password: "magic-password"}
	return
}

// GetUser
func (c *UserController) DeleteUser(at struct{
	at.DeleteMapping `value:"/{id}"`
	at.Operation   `id:"Update Employee" description:"Delete User"`
}, id int) (response *UserResponse) {
	response = new(UserResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("Success")
	return
}




