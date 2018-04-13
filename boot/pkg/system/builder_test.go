package system

import (
	"testing"
	"github.com/hidevopsio/hi/boot/pkg/utils"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hi/boot/pkg/log"
)

const (
	application = "application"
	config      = "/config"
	yaml        = "yaml"
)

func TestBuilder(t *testing.T) {
	b := &Builder{
		Path:       utils.GetWorkingDir("boot/pkg/system/builder_test.go") + config,
		Name:       application,
		FileType:   yaml,
		Profile:    "local",
		ConfigType: Configuration{},
	}

	cp, err := b.Build()
	assert.Equal(t, nil, err)
	c := cp.(*Configuration)
	log.Print(c)
}
