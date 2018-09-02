package cmd

import (
	"github.com/hidevopsio/hiboot/pkg/app/cli"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFooCommands(t *testing.T) {
	testApp := cli.NewTestApplication(t, new(fooCommand))

	t.Run("should run foo command", func(t *testing.T) {
		_, err := testApp.RunTest()
		assert.Equal(t, nil, err)
	})
}
