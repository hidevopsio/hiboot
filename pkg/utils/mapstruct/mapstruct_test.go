package mapstruct

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type Foo struct {
	Name string
}

func TestDecode(t *testing.T) {
	var foo Foo

	src := map[string]string {
		"name": "foo",
	}

	t.Run("should decode map to struct", func(t *testing.T) {
		err := Decode(&foo, src)
		assert.Equal(t, nil, err)
		assert.Equal(t, "foo", foo.Name)
	})

	t.Run("should return result must be a pointer", func(t *testing.T) {
		err := Decode(nil, src)
		assert.Equal(t, "result must be a pointer", err.Error())
	})

	t.Run("should return parameters should not be nil", func(t *testing.T) {
		err := Decode(&foo, nil)
		assert.Equal(t, "parameters should not be nil", err.Error())
	})
}
