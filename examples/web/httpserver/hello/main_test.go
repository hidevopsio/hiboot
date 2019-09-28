package hello

import (
	"hidevops.io/hiboot/pkg/app/web"
	"net/http"
	"testing"
	"time"
)

func TestRunMain(t *testing.T) {
	go main()
}

func TestController(t *testing.T) {
	time.Sleep(time.Second)
	web.NewTestApp().
		SetProperty(web.ViewEnabled, true).
		Run(t).
		Get("/").
		Expect().Status(http.StatusOK)
}
