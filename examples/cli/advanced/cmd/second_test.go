package cmd

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/starter/cli"
	"github.com/stretchr/testify/assert"
)

func TestSecondCommands(t *testing.T) {
	testApp := cli.NewTestApplication(new(secondCommand))

	t.Run("should run second command", func(t *testing.T) {
		_, err := testApp.RunTest("")
		assert.Equal(t, nil, err)
	})
}

