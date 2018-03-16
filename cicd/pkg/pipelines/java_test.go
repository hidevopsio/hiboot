package pipelines

import (
	"testing"
	"github.com/hidevopsio/hi/boot/pkg/log"
)

func init()  {
	log.SetLevel("debug")
}

func TestPipeline(t *testing.T)  {
	p := &JavaPipeline{
		Pipeline{
			App: "test",
			Project: "demo",
		},
	}
	p.EnsurePipeline()
	p.PullSourceCode()
	p.Deploy()
}