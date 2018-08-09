package inject

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
)

func TestTag(t *testing.T) {
	tag := new(BaseTag)

	t.Run("should get properties", func(t *testing.T) {
		p := tag.Properties()
		assert.NotEqual(t, nil, p)
	})

	t.Run("should get check if it's singleton", func(t *testing.T) {
		s := tag.IsSingleton()
		assert.Equal(t, false, s)
	})

	t.Run("should get properties", func(t *testing.T) {
		fakeObj := struct{Name string}{}
		objVal := reflect.ValueOf(fakeObj)
		field := objVal.Type().Field(0)
		f := tag.Decode(objVal, field, "fake")
		assert.Equal(t, nil, f)
	})
}