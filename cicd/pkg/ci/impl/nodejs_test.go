package impl

import (
	"testing"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"os"
	"github.com/hidevopsio/hi/cicd/pkg/ci"
	"github.com/magiconair/properties/assert"
)

func TestNodeJsPipeline(t *testing.T) {

	log.Debug("Test NodeJs Pipeline")

	nodeJs := &NodeJsPipeline{}
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	pi := &ci.Pipeline{
		Name:    "nodejs",
		Project: "demo",
		Profile: "dev",
		App:     "hello-angular",
		Scm:     ci.Scm{Url: os.Getenv("SCM_URL")},
	}
	nodeJs.Init(pi)
	err := nodeJs.Run(username, password, false)
	assert.Equal(t, nil, err)
}
