package k8s

import (
	"github.com/hidevopsio/hi/boot/pkg/log"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hi/cicd/pkg/orch"
)

func init()  {
	log.SetLevel(log.DebugLevel)
}

func TestServiceCreation(t *testing.T) {
	log.Debug("TestServiceCreation()")

	projectName := "demo"
	profile     := "dev"
	namespace   := projectName + "-" + profile
	app         := "hello-world"

	p := []orch.Ports{
		{
			Name: "8080-tcp",
			Port: 8080,
		},
		{
			Name: "7575-tcp",
			Port: 7575,
		},
	}

	service := NewService(app, namespace)
	err := service.Create(p)
	assert.Equal(t, nil, err)
}


func TestServiceDeletion(t *testing.T) {
	log.Debug("TestServiceDeletion()")

	projectName := "demo"
	profile     := "dev"
	namespace   := projectName + "-" + profile
	app         := "hello-world"

	service := NewService(app, namespace)
	err := service.Delete()
	assert.Equal(t, nil, err)
}

