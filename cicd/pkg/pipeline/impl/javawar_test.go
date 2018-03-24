package impl

import (
	"testing"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"github.com/hidevopsio/hi/cicd/pkg/pipeline"
	"github.com/hidevopsio/hi/cicd/pkg/ci"
)

func init()  {
	log.SetLevel("debug")
}

func TestJavaWarPipeline(t *testing.T)  {

	log.Debug("Test Java War Pipeline")

	javaWarPipeline := &JavaWarPipeline{
		JavaPipeline{
			pipeline.Pipeline{
				App: "test",
				Project: "demo",
			},
		},
	}

	ci.Run(javaWarPipeline)
}
