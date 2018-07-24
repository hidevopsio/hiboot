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

	boltRepo := configuration.NewRepository("foo")
	bolt := boltRepo.DataSource().(*bolt)
	assert.NotEqual(t, nil, bolt)
	bolt.Close()
}
