package pipeline


import (
	"testing"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"github.com/stretchr/testify/assert"
)

func init()  {
	log.SetLevel("debug")
}

func TestPipelineBuilder(t *testing.T)  {

	log.Debug("Test Pipeline Builder")

	cfg := Build("java")

	log.Debug(cfg)

	assert.Equal(t, "java", cfg.Pipeline.Name)
}
