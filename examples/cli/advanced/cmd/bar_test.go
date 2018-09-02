package cmd

import (
	"github.com/hidevopsio/hiboot/pkg/app/cli"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBarCommands(t *testing.T) {

	testApp := cli.NewTestApplication(t, new(barCommand))

	t.Run("should run bar command", func(t *testing.T) {
		_, err := testApp.RunTest()
		assert.Equal(t, nil, err)
	})

	t.Run("should run baz command", func(t *testing.T) {
		_, err := testApp.RunTest("baz")
		assert.Equal(t, nil, err)
	})
}
