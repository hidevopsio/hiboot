package main

import (
	"hidevops.io/hiboot/pkg/app/web"
	"net/http"
	"testing"
)

func TestRunMain(t *testing.T) {
	go main()
}

func TestController(t *testing.T) {
	testApp := web.NewTestApp(t).Run(t)

	t.Run("should get /static/ui/index.html ", func(t *testing.T) {
		testApp.Get("/static/ui").
			Expect().Status(http.StatusOK).
			Body().Contains("Hiboot Web Application Example")
	})

	t.Run("should get /static/simple/ui/index.html ", func(t *testing.T) {
		testApp.Get("/static/simple/ui").
			Expect().Status(http.StatusOK).
			Body().Contains("Hiboot Web Application Example")
	})

}
