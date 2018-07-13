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

	err := Decode(&foo, src)
	assert.Equal(t, nil, err)
	assert.Equal(t, "foo", foo.Name)

	err = Decode(nil, src)
	assert.Equal(t, "result must be a pointer", err.Error())

}
