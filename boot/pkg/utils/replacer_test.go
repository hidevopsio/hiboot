package utils

import (
	"testing"
	"os"
	"regexp"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"github.com/stretchr/testify/assert"
)

func ParseVariables(src string) [][]string    {
	matches := regexp.MustCompile(`\$\{(.*?)\}`).FindAllStringSubmatch(src, -1)
	if matches == nil {
		log.Println("No matches found.")
		return nil
	}
	return matches
}

func TestParseVariables(t *testing.T) {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "bar")
	source := "the-${FOO}-${BAR}-env"

	matches := ParseVariables(source)

	log.Print(matches)
	assert.Equal(t, "${FOO}", matches[0][0])
	assert.Equal(t, "FOO", matches[0][1])
	assert.Equal(t, "${BAR}", matches[1][0])
	assert.Equal(t, "BAR", matches[1][1])
}