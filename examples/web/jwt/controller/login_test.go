package controller

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"net/http"
	"testing"
)

func TestFooLogin(t *testing.T) {

	testApp := web.NewTestApplication(t, new(loginController))

	testApp.Post("/login").
		WithJSON(userRequest{Username: "mike", Password: "daDg83t"}).
		Expect().Status(http.StatusOK)
}
