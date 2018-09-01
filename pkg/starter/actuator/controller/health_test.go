package controller_test

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
	_ "github.com/hidevopsio/hiboot/pkg/starter/actuator/controller"
	"net/http"
	"testing"
)

func TestHealthController(t *testing.T) {
	web.NewTestApplication(t).
		Get("/health").
		Expect().Status(http.StatusOK)
}
