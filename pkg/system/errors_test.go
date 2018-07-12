package system

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInvalidControllerError(t *testing.T) {
	err := InvalidControllerError{Name: "TestController"}

	assert.Equal(t, "TestController must be derived from web.Controller", err.Error())
}

func TestNotFoundError(t *testing.T) {
	err := NotFoundError{Name: "TestObject"}

	assert.Equal(t, "TestObject is not found", err.Error())
}
