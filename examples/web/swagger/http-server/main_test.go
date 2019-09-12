package main

import (
	"hidevops.io/hiboot/pkg/app/web"
	"net/http"
	"testing"
	"time"
)

func TestController(t *testing.T) {
	testApp := web.NewTestApp(t).Run(t)

	t.Run("should get index.html ", func(t *testing.T) {
		testApp.Get("/public/ui").
			Expect().Status(http.StatusOK)
	})

	t.Run("should get hello.txt ", func(t *testing.T) {
		testApp.Get("/public/ui/hello.txt").
			Expect().Status(http.StatusOK)
	})
}

func TestRunMain(t *testing.T) {
	go main()
	time.Sleep(100 * time.Millisecond)
}