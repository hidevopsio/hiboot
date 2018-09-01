package cli

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	rootCmd := new(rootCommand)

	rootCmd.Init()
	err := rootCmd.Run(nil)
	assert.Equal(t, nil, err)
}
