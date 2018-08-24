package cmd

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/app/cli"
	"github.com/stretchr/testify/assert"
)

func TestFirstCommands(t *testing.T) {
	testApp := cli.NewTestApplication(new(FirstCommand))

	t.Run("should run first command", func(t *testing.T) {
		_, err := testApp.RunTest("-t", "10")
		assert.Equal(t, nil, err)
	})

	t.Run("should run second command", func(t *testing.T) {
		_, err := testApp.RunTest("second")
		assert.Equal(t, nil, err)
	})

	t.Run("should run foo command", func(t *testing.T) {
		_, err := testApp.RunTest("second", "foo")
		assert.Equal(t, nil, err)
	})

	t.Run("should run bar command", func(t *testing.T) {
		_, err := testApp.RunTest("second", "bar")
		assert.Equal(t, nil, err)
	})

	t.Run("should report unknown command", func(t *testing.T) {
		_, err := testApp.RunTest("not-exist-command")
		assert.Contains(t, err.Error(),"unknown command")
	})
}
