package starter

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"path/filepath"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type FakeProperties struct {
	Name string
}

type FakeConfiguration struct {
	FakeProperties FakeProperties `mapstructure:"fake"`
}

type Foo struct {
	Name string
}

func init() {
	log.SetLevel(log.DebugLevel)
	utils.EnsureWorkDir("../../")
	Add("fake", FakeConfiguration{})
}

func (c *FakeConfiguration) Foo() *Foo {
	f := new(Foo)
	f.Name = c.FakeProperties.Name

	return f
}

func TestBuild(t *testing.T) {
	configPath := filepath.Join(utils.GetWorkDir(), "config")
	fakeFile := "application-fake.yaml"

	fakeContent :=
		"fake:\n" +
		"  name: foo\n"
	n, err := utils.WriterFile(configPath, fakeFile, []byte(fakeContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, n, len(fakeContent))

	config := GetInstance()
	config.Build()

	fc := config.Configuration("fake").(*FakeConfiguration)
	assert.Equal(t, "foo", fc.FakeProperties.Name)
	assert.Equal(t, "foo", config.Instances()["foo"].(*Foo).Name)

	os.Remove(filepath.Join(configPath, fakeFile))
}