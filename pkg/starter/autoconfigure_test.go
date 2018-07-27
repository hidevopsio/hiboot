package starter

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"path/filepath"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"os"
)

type FakeProperties struct {
	Name string
	Nickname string
	Username string
}

type FakeConfiguration struct {
	Fake FakeProperties `mapstructure:"fake"`
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
	f.Name = c.Fake.Name

	return f
}

func TestBuild(t *testing.T) {
	configPath := filepath.Join(utils.GetWorkDir(), "config")
	fakeFile := "application-fake.yaml"
	os.Remove(filepath.Join(configPath, fakeFile))
	fakeContent :=
		"fake:\n" +
		"  name: foo\n" +
		"  nickname: ${app.name} ${fake.name}\n" +
		"  username: ${unknown.name:bar}\n"
	n, err := utils.WriterFile(configPath, fakeFile, []byte(fakeContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, n, len(fakeContent))

	config := GetAutoConfiguration()
	config.Build()
	fci := config.Configuration("fake")
	assert.NotEqual(t, nil, fci)
	fc := fci.(*FakeConfiguration)

	assert.Equal(t, "hiboot foo", fc.Fake.Nickname)
	assert.Equal(t, "bar", fc.Fake.Username)
	assert.Equal(t, "foo", fc.Fake.Name)
	assert.Equal(t, "foo", config.Instances()["foo"].(*Foo).Name)
}