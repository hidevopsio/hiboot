package bolt

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBoltCrd(t *testing.T) {
	testBucket := "test-bucket"
	testKey := "hello"
	testValue := "world"
	
	dataSource := make(map[string]interface{})

	dataSource["database"] = "test.db"
	dataSource["mode"] = 0600
	dataSource["timeout"] = 2

	b := &Bolt{}

	t.Run("should open bolt database", func(t *testing.T) {
		err := b.Open(dataSource)
		assert.Equal(t, nil, err)
	})

	b.SetNamespace(testBucket)

	t.Run("should put data into bolt database", func(t *testing.T) {
		err := b.Put([]byte(testKey), []byte(testValue))
		assert.Equal(t, nil, err)
	})

	t.Run("should get data into bolt database", func(t *testing.T) {
		val, err := b.Get([]byte(testKey))
		assert.Equal(t, nil, err)
		assert.Equal(t, testValue, string(val))
	})

	t.Run("should delete data into bolt database", func(t *testing.T) {
		err := b.Delete([]byte(testKey))
		assert.Equal(t, nil, err)
	})
	// close bolt database
	b.Close()
}

func TestBoltErrorHandling(t *testing.T) {

	b := &Bolt{}

	dataSource := make(map[string]interface{})

	dataSource["database"] = "test.db"
	dataSource["mode"] = 0600
	dataSource["timeout"] = 2

	t.Run("should return error if dataSource input is nil ", func(t *testing.T) {
		err := b.Open(nil)
		assert.Equal(t, "parameters of mapstruct.Decode must not be nil", err.Error())
	})

	t.Run("should return error ", func(t *testing.T) {
		err := b.Open(dataSource)
		assert.Equal(t, nil, err)
		b.Close()
	})

}