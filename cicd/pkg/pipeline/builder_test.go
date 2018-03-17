package pipeline


import (
	"testing"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hi/boot/pkg/system"
)

func init()  {
	log.SetLevel("debug")
}

func TestPipelineBuilder(t *testing.T)  {

	log.Debug("Test Pipeline Builder")

	syscfg := system.Build()
	log.Debug(syscfg)
	assert.Equal(t, "hi", syscfg.App.Name)

	cfg := Build("java")
	log.Debug(cfg)
	assert.Equal(t, "java", cfg.Pipeline.Name)
}
