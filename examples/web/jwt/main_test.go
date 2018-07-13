package main

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"net/http"
	"github.com/hidevopsio/hiboot/examples/web/jwt/controllers"
	"os"
)


func TestController(t *testing.T) {
	web.NewTestApplication(t).
		Get("/foo").
		WithQueryObject(controllers.FooRequest{Name: "Peter", Age: 18}).
		Expect().Status(http.StatusOK)
}

func TestMain(m *testing.M) {
	go func() {
		os.Exit(m.Run())
	}()
}