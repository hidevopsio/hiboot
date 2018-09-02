package cmd

import (
	"github.com/hidevopsio/hiboot/pkg/app/cli"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecondCommands(t *testing.T) {
	testApp := cli.NewTestApplication(t, new(secondCommand))

	t.Run("should run second command", func(t *testing.T) {
		_, err := testApp.RunTest("")
		assert.Equal(t, nil, err)
	})
}
