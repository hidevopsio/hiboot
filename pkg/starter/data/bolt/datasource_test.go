package bolt

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDataSource(t *testing.T) {

	properties := &properties{
		Database: "test.db",
		Mode: 0600,
		Timeout: 2,
	}

	d := GetDataSource()

	t.Run("should open bolt database", func(t *testing.T) {
		err := d.Open(nil)
		assert.Equal(t, InvalidPropertiesError, err)
	})

	t.Run("should open bolt database", func(t *testing.T) {
		err := d.Open(properties)
		assert.Equal(t, nil, err)
	})

	// close bolt database
	d.Close()
}


func TestDataSourceWithEmptyFile(t *testing.T) {

	properties := &properties{
		Timeout: 2,
	}
	d := GetDataSource()
	t.Run("should open bolt database", func(t *testing.T) {
		err := d.Open(properties)
		assert.NotEqual(t, nil, err)
	})
}
