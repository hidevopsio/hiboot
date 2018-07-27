package bolt

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewBolt(t *testing.T) {
	configuration := new(configuration)
	configuration.BoltProperties = properties{
		Database: "test.db",
		Mode: 0600,
		Timeout: 1,
	}

	repository := configuration.BoltRepository()
	assert.NotEqual(t, nil, repository)
	repository.DataSource().(DataSource).Close()
}

func TestNewBoltWithError(t *testing.T) {
	configuration := new(configuration)

	repository := configuration.BoltRepository()
	assert.NotEqual(t, nil, repository)

	configuration.BoltProperties = properties{
		Timeout: 1,
	}

	repository = configuration.BoltRepository()
	assert.NotEqual(t, nil, repository)
}