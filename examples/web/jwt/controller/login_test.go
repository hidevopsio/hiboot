package controller

import (
	"testing"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/app/web"
)

func TestFooLogin(t *testing.T) {

	testApp := web.NewTestApplication(t, new(loginController))

	testApp.Post("/login").
		WithJSON(userRequest{Username: "mike", Password: "daDg83t"}).
		Expect().Status(http.StatusOK)
}
