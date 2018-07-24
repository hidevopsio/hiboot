package bolt

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBoltCrd(t *testing.T) {
	testBucket := []byte("test-bucket")
	testKey := []byte("hello")
	testValue := []byte("world")

	properties := &properties{
		Database: "test.db",
		Mode: 0600,
		Timeout: 1,
	}

	b := GetInstance()

	t.Run("should open bolt database", func(t *testing.T) {
		err := b.Open(properties)
		assert.Equal(t, nil, err)
	})

	t.Run("should put data into bolt database", func(t *testing.T) {
		err := b.Put(testBucket, testKey, testValue)
		assert.Equal(t, nil, err)
	})

	t.Run("should get data into bolt database", func(t *testing.T) {
		val, err := b.Get(testBucket, testKey)
		assert.Equal(t, nil, err)
		assert.Equal(t, testValue, val)
	})

	t.Run("should delete data into bolt database", func(t *testing.T) {
		err := b.Delete(testBucket, testKey)
		assert.Equal(t, nil, err)
	})
	// close bolt database
	b.Close()
}

func TestBoltErrorHandling(t *testing.T) {

	b := GetInstance()

	properties := &properties{
		Database: "test.db",
		Mode: 0600,
		Timeout: 2,
	}

	t.Run("should return error if dataSource input is nil ", func(t *testing.T) {
		err := b.Open(nil)
		assert.Equal(t, "properties must not be nil", err.Error())
	})

	t.Run("should return error ", func(t *testing.T) {

		err := b.Open(properties)
		assert.Equal(t, nil, err)
		b.Close()
	})

	t.Run("should return error ", func(t *testing.T) {

		err := b.Open(properties)
		assert.Equal(t, nil, err)
		err = b.Open(properties)
		assert.Equal(t, nil, err)
		b.Close()
	})

}