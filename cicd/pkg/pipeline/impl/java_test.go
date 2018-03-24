package impl

import (
	"testing"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"github.com/hidevopsio/hi/cicd/pkg/pipeline"
)

func init()  {
	log.SetLevel("debug")
}

func run(p pipeline.PipelineInterface)  {
	p.EnsureParam()
	p.Build()
	p.Deploy()
}

func TestJavaPipeline(t *testing.T)  {

	log.Debug("Test Java Pipeline")

	javaPipeline := &JavaPipeline{
		pipeline.Pipeline{
			App: "test",
			Project: "demo",
		},
	}

	run(javaPipeline)
}
