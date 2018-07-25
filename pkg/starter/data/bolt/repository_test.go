package bolt

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var (
	testBucket = "test-bucket"
	testKey = []byte("hello")
	testValue = []byte("world")
)


func TestRepositoryCrd(t *testing.T) {

	properties := &properties{
		Database: "test.db",
		Mode: 0600,
		Timeout: 2,
	}

	b := GetInstance()

	t.Run("should open bolt database", func(t *testing.T) {
		err := b.Open(properties)
		assert.Equal(t, nil, err)
	})

	r := &repository{}
	r.SetDataSource(b)

	r.SetName(testBucket)
	assert.Equal(t, testBucket, r.Name())

	t.Run("should put data into bolt database", func(t *testing.T) {
		err := r.Put(testKey, testValue)
		assert.Equal(t, nil, err)
	})

	t.Run("should get data into bolt database", func(t *testing.T) {
		val, err := r.Get(testKey)
		assert.Equal(t, nil, err)
		assert.Equal(t, testValue, val)
	})

	t.Run("should delete data into bolt database", func(t *testing.T) {
		err := r.Delete(testKey)
		assert.Equal(t, nil, err)
	})

	// close bolt database
	r.DataSource().(*bolt).Close()
}

func TestRepositoryWithNilDataSource(t *testing.T) {
	r := &repository{}

	t.Run("should put data into bolt database", func(t *testing.T) {
		err := r.Put(testKey, testValue)
		assert.Equal(t, "dataSource is nil", err.Error())
	})

	t.Run("should get data into bolt database", func(t *testing.T) {
		_, err := r.Get(testKey)
		assert.Equal(t, "dataSource is nil", err.Error())
	})

	t.Run("should delete data into bolt database", func(t *testing.T) {
		err := r.Delete(testKey)
		assert.Equal(t, "dataSource is nil", err.Error())
	})
}
