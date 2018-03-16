package pipelines

import (
	"testing"
	"github.com/hidevopsio/hi/boot/pkg/log"
)

func init()  {
	log.SetLevel("debug")
}

func TestPipeline(t *testing.T)  {

	log.Debug("Test Java Pipeline")

	javaPipeline := &JavaPipeline{
		Pipeline{
			App: "test",
			Project: "demo",
		},
	}

	Run(javaPipeline)
}
