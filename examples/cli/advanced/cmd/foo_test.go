package cmd

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/starter/cli"
	"github.com/stretchr/testify/assert"
)

func TestFooCommands(t *testing.T) {
	testApp := cli.NewTestApplication(new(fooCommand))

	t.Run("should run foo command", func(t *testing.T) {
		_, err := testApp.RunTest()
		assert.Equal(t, nil, err)
	})
}

