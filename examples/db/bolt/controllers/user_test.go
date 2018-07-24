package controllers

import (
	"testing"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/examples/db/bolt/domain"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"github.com/hidevopsio/hiboot/pkg/log"
)


func init() {
	log.SetLevel(log.DebugLevel)
	utils.EnsureWorkDir("..")
}

func TestCrdRequest(t *testing.T) {
	app := web.NewTestApplication(t, new(UserController))

	t.Run("should add user with POST request", func(t *testing.T) {
		// First, let's Post User
		app.Post("/user").
			WithJSON(domain.User{Id: "1", Name: "Peter", Age: 18}).
			Expect().Status(http.StatusOK)
	})

	t.Run("should get user with GET request", func(t *testing.T) {
		// Then Get User
		app.Get("/user").
			WithQuery("id", "1").
			Expect().Status(http.StatusOK)
	})

	t.Run("should return 404 if trying to find a record that does not exist", func(t *testing.T) {
		// Then Get User
		app.Get("/user").
			WithQuery("id", "9999").
			Expect().Status(http.StatusNotFound)
	})

	t.Run("should delete the record with DELETE request", func(t *testing.T) {
		// Finally Delete User
		app.Delete("/user").
			WithQuery("id", "1").
			Expect().Status(http.StatusOK)
	})

}
