package controllers

import (
	"testing"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/examples/db/bolt/domain"
	"github.com/hidevopsio/hiboot/pkg/utils"
)


func init() {
	utils.ChangeWorkDir("../")
}

func TestCrdRequest(t *testing.T) {
	app := web.NewTestApplication(t, new(UserController))

	// First, let's Post User
	app.Post("/user").
		WithJSON(domain.User{Id: "1", Name: "Peter", Age: 18}).
		Expect().Status(http.StatusOK)

	// Then Get User
	app.Get("/user").
		WithQuery("id", "1").
		Expect().Status(http.StatusOK)

	// Then Get User
	app.Get("/user").
		WithQuery("id", "9999").
		Expect().Status(http.StatusNotFound)

	// Finally Delete User
	app.Delete("/user").
		WithQuery("id", "1").
		Expect().Status(http.StatusOK)

}
