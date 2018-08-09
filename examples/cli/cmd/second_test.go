package cmd

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/starter/cli"
	"github.com/stretchr/testify/assert"
)

func TestSecondCommands(t *testing.T) {
	secondCmd := new(secondCommand)
	firstCmd := new(firstCommand)
	firstCmd.Add(secondCmd)
	testApp := cli.NewTestApplication(firstCmd)

	t.Run("should run second command", func(t *testing.T) {
		_, err := testApp.RunTest("first", "second")
		assert.Equal(t, nil, err)
	})
}
