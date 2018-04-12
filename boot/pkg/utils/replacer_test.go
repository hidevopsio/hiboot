package utils

import (
	"testing"
	"github.com/magiconair/properties/assert"
)

func TestEnvReplace(t *testing.T) {

	env := "asd-${_JAVA_OPTIONS}-qwe"

	NewEnv := EnvReplace(env)
	assert.Equal(t, "asd--Djava.net.preferIPv4Stack=true-qwe", NewEnv)
}