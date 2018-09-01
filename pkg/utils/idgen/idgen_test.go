package idgen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNext(t *testing.T) {

	t.Run("should generate id in uint", func(t *testing.T) {
		id, err := Next()
		assert.Equal(t, nil, err)
		assert.NotEqual(t, 0, id)
	})

	t.Run("should generate id in string", func(t *testing.T) {
		id, err := NextString()
		assert.Equal(t, nil, err)
		assert.NotEqual(t, 0, id)
	})

}
