package system

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"github.com/magiconair/properties/assert"
)

const (
	application = "application"
	config      = "/config"
	yaml        = "yaml"
)

func TestBuilder(t *testing.T) {
	b := &Builder{
		Path:       utils.GetWorkingDir("") + config,
		Name:       application,
		FileType:   yaml,
		Profile:    "local",
		ConfigType: Configuration{},
	}

	cp, err := b.Build()
	assert.Equal(t, nil, err)

	c := cp.(*Configuration)
	assert.Equal(t, "hiboot", c.App.Name)
}
