package cmd

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/starter/cli"
	"github.com/stretchr/testify/assert"
)

func TestFooCommands(t *testing.T) {
	fooCmd := new(fooCommand)
	secondCmd := new(secondCommand)
	firstCmd := new(firstCommand)
	secondCmd.Add(fooCmd)
	firstCmd.Add(secondCmd)
	testApp := cli.NewTestApplication(firstCmd)

	t.Run("should run second command", func(t *testing.T) {
		_, err := testApp.RunTest("first", "second", "foo")
		assert.Equal(t, nil, err)
	})
}

