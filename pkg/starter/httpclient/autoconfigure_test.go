package httpclient

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestConfiguration(t *testing.T) {
	c := newConfiguration()

	t.Run("should get a struct", func(t *testing.T) {
		client := c.Client()
		assert.IsType(t, reflect.Struct, reflect.TypeOf(client).Kind())
	})

}
