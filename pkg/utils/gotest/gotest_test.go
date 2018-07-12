package gotest

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestIsRunning(t *testing.T) {
	isTestRunning := IsRunning()

	assert.Equal(t, true, isTestRunning)
}