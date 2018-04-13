package utils

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"regexp"
	"os"
	"github.com/hidevopsio/hi/boot/pkg/log"
)

type Bar struct {
	Name string
	Profile string
	SubBar SubBar
}

type Foo struct{
	Name string
	Project string
	Bar Bar
}

type SubBar struct {
	Name string
}

func TestReplaceVariable(t *testing.T)  {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "bar")
	f := &Foo{
		Name: "foo",
		Project: "it's ${FOO} project",
		Bar: Bar{
			Name: "my name is ${BAR}",
			Profile: "${name}-bar",
			SubBar: SubBar{
				Name: "${bar.name}",
			},
		},
	}

	err := Replace(f, f)
	assert.Equal(t, nil, err)
	assert.Equal(t, "foo-bar", f.Bar.Profile)

}

func TestParseVariables(t *testing.T) {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "bar")
	os.Setenv("foo.bar", "fb")
	source := "the-${FOO}-${BAR}-${foo.bar}-env"

	re := regexp.MustCompile(`\$\{(.*?)\}`)

	matches := ParseVariables(source, re)

	log.Print(matches)
	assert.Equal(t, "${FOO}", matches[0][0])
	assert.Equal(t, "FOO", matches[0][1])
	assert.Equal(t, "${BAR}", matches[1][0])
	assert.Equal(t, "BAR", matches[1][1])
}

func TestReplaceReferences(t *testing.T)  {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "bar")
	f := &Foo{
		Name: "foo",
		Project: "it's ${FOO} project",
		Bar: Bar{
			Name: "my name is ${BAR}",
			Profile: "${name}-bar",
			SubBar: SubBar{
				Name: "${bar.name}",
			},
		},
	}
	res, err := ParseReferences(f, []string{"name"})
	assert.Equal(t, nil, err)
	log.Println("res: ", res)

	res, err = ParseReferences(f, []string{"bar", "name"})
	assert.Equal(t, nil, err)
	log.Println("res: ", res)
}