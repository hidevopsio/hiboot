package main

import (
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/app/web/server"
	"net/http"
	"testing"
)

func TestRunMain(t *testing.T) {
	go main()
}

func TestController(t *testing.T) {
	testApp := web.NewTestApp(t).SetProperty(server.ContextPath, basePath).Run(t)

	t.Run("should get employee ", func(t *testing.T) {
		testApp.Get("/api/v1/greeting-server/hello").
			Expect().Status(http.StatusOK)
	})

	t.Run("should get employee ", func(t *testing.T) {
		testApp.Get("/api/v1/greeting-server/hey").
			Expect().Status(http.StatusOK)
	})

}

