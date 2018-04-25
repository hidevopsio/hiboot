package system

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"github.com/magiconair/properties/assert"
	"path/filepath"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
)


func TestBuilderBuild(t *testing.T) {

	b := &Builder{
		Path:       filepath.Join(utils.GetWorkingDir(""), "config"),
		Name:       "application",
		FileType:   "yaml",
		Profile:    "local",
		ConfigType: Configuration{},
	}

	cp, err := b.Build()
	assert.Equal(t, nil, err)

	c := cp.(*Configuration)
	assert.Equal(t, "hiboot", c.App.Name)

	log.Print(c)
}


func TestBuilderInit(t *testing.T) {
	b := &Builder{
		Path:       filepath.Join(os.TempDir(), "config"),
		Name:       "foo",
		FileType:   "yaml",
		ConfigType: Configuration{},
	}

	err := b.Init()
	assert.Equal(t, nil, err)
}

func TestBuilderSave(t *testing.T) {
	b := &Builder{
		Path:       filepath.Join(os.TempDir(), "config"),
		Name:       "foo",
		FileType:   "yaml",
		ConfigType: Configuration{},
	}

	err := b.Init()
	assert.Equal(t, nil, err)

	c := &Configuration{
		App: App{
			Name: "foo",
			Project: "bar",
		},
		Server: Server{
			Port: 8080,
		},
	}
	err = b.Save(c)
	assert.Equal(t, nil, err)
}
