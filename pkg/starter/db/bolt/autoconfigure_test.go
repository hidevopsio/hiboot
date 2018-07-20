package bolt

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewBolt(t *testing.T) {
	configuration := new(Configuration)
	configuration.BoltProperties = Properties{
		Database: "test.db",
		Mode: 0600,
		Timeout: 1,
	}

	boltRepo := configuration.NewRepository()
	bolt := boltRepo.DataSource().(*Bolt)
	assert.NotEqual(t, (*Bolt)(nil), bolt)
	bolt.Close()
}
