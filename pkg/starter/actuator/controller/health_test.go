package controller_test

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"net/http"
	_ "github.com/hidevopsio/hiboot/pkg/starter/actuator/controller"
)

func TestHealthController(t *testing.T) {
	web.NewTestApplication(t).
		Get("/health").
		Expect().Status(http.StatusOK)
}
