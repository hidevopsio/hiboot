package cmd

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/starter/cli"
	"github.com/stretchr/testify/assert"
)

func TestFirstCommands(t *testing.T) {
	testApp := cli.NewTestApplication(new(firstCommand))

	t.Run("should run first command", func(t *testing.T) {
		_, err := testApp.RunTest("first")
		assert.Equal(t, nil, err)
	})
}
