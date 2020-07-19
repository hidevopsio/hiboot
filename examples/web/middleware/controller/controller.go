package controller

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
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

type UserRequests struct {
	at.RequestParams
	at.Schema

	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty" json:"page,omitempty" validate:"min=1"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty" json:"per_page,omitempty" validate:"min=1"`

	Expr string `json:"expr"`
}

type UserResponse struct {
	at.ResponseBody `json:"-"`
	model.BaseResponse
	Data *User `json:"data"`
}

type ListUserResponse struct {
	at.ResponseBody `json:"-"`
	model.BaseResponse
	Data *UserRequests `json:"data"`
}

// GetUser
func (c *UserController) GetUser(_ struct{
	at.GetMapping `value:"/{id}"`
	at.Operation   `id:"Update Employee" description:"Get User by ID"`
	// /user/{id} -> `values:"user:read" type:"path" in:"id"`
	at.RequiresPermissions `values:"user:read" type:"path" in:"id"`
}, id int, ctx context.Context) (response *UserResponse) {
	response = new(UserResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("Success")
	response.Data = &User{ID: id, Username: "john.deng", Password: "magic-password"}
	return
}

// GetUser
func (c *UserController) GetUserQuery(_ struct{
	at.GetMapping `value:"/query"`
	at.Operation   `id:"Update Employee" description:"Query User"`
	// /user?id=12345 -> `values:"user:read" type:"query" in:"id"`
	at.RequiresPermissions `values:"user:read" type:"query" in:"id"`
}, ctx context.Context) (response *UserResponse) {
	response = new(UserResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("Success")
	id, _ := ctx.URLParamInt("id")
	response.Data = &User{ID: id, Username: "john.deng", Password: "magic-password"}
	return
}


// GetUser
func (c *UserController) GetUsers(_ struct{
	at.GetMapping `value:"/"`
	at.Operation   `id:"Update Employee" description:"Get User List"`
	at.RequiresPermissions `values:"user:list,team:*" type:"query:pagination" in:"page,per_page" out:"expr"`
}, request *UserRequests, ctx context.Context) (response *ListUserResponse) {
	log.Debugf("expr: %v", request.Expr)
	log.Debugf("header.expr: %v", ctx.GetHeader("expr"))
	response = new(ListUserResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("Success")
	response.Data = request
	return
}


// GetUser
func (c *UserController) DeleteUser(_ struct{
	at.DeleteMapping `value:"/{id}"`
	at.Operation   `id:"Update Employee" description:"Delete User by ID"`
}, id int) (response *UserResponse) {
	response = new(UserResponse)
	response.SetCode(http.StatusOK)
	response.SetMessage("Success")
	return
}




