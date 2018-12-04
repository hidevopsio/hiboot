package system

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadYamlProperties(t *testing.T) {

	res, err := ReadYamlFromFile("config/test-file.yml")

	assert.Equal(t, nil, err)
	assert.Equal(t, "foo name", res["foo"])

	res, err = ReadYamlFromFile("config/application.yml")

	assert.Equal(t, nil, err)
	assert.Equal(t, "foo name", res["foo"])
}
