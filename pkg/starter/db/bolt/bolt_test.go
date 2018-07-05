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

	err := b.Open(dataSource)
	assert.Equal(t, nil, err)

	b.SetNamespace(testBucket)

	err = b.Put([]byte(testKey), []byte(testValue))
	assert.Equal(t, nil, err)

	val, err := b.Get([]byte(testKey))
	assert.Equal(t, nil, err)
	assert.Equal(t, testValue, string(val))

	err = b.Delete([]byte(testKey))
	assert.Equal(t, nil, err)

}