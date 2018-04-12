package impl

import (
	"testing"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"os"
	"github.com/hidevopsio/hi/cicd/pkg/ci"
	"github.com/magiconair/properties/assert"
)

/*func TestPipelineInit(t *testing.T)  {
	p := &ci.Pipeline{
		Name: "nodejs",
		Profile: "dev",
		App: "hello-world",
		Project: "demo",
		Scm: ci.Scm{Url: os.Getenv("SCM_URL")},
	}
	pipelineFactory := new(factories.PipelineFactory)
	pipeline, err := pipelineFactory.New(p.Name)
	if err !=nil {

	}

	log.Info(pipeline)
	//assert.Equal(t,&ci.Pipeline{} , pipeline)
}*/

func TestNodeJsPipeline(t *testing.T)  {

	log.Debug("Test NodeJs Pipeline")

	nodeJs := &NodeJsPipeline{}
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	pi :=&ci.Pipeline{
		Name: "nodejs",
		Profile: "dev",
		App: "admin",
		Project: "demo",
		Scm: ci.Scm{Url: os.Getenv("SCM_URL")},
	}
	nodeJs.Init(pi)
	err := nodeJs.Run(username, password, false)
	assert.Equal(t, nil, err)
}