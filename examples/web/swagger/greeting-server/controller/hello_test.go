package controller

import (
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/app/web/server"
	"net/http"
	"sync"
	"testing"
)

var mu sync.Mutex
func TestController(t *testing.T) {
	mu.Lock()
	basePath := "/api/v1/greeting-server"
	testApp := web.NewTestApp(t, newHelloController).SetProperty(server.ContextPath, basePath).Run(t)

	t.Run("should get employee ", func(t *testing.T) {
		testApp.Get(basePath + "/hello").
			Expect().Status(http.StatusOK)
	})

	t.Run("should get employee ", func(t *testing.T) {
		testApp.Get(basePath + "/hey").
			Expect().Status(http.StatusOK)
	})
	mu.Unlock()
}

