package cli

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	rootCmd := new(rootCommand)

	rootCmd.Init()
	err := rootCmd.Run(nil)
	assert.Equal(t, nil, err)
}
