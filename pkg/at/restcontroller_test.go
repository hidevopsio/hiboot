package at

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestControllerAnnotations(t *testing.T) {

	t.Run("should pass test for annotation RequestMapping", func(t *testing.T) {
		rm := RequestMapping{HttpMethod{Annotation{"/foo"}}}

		assert.Equal(t, rm.Annotation.Value, "/foo")
		assert.Equal(t, rm, RequestMapping{HttpMethod{Annotation{"/foo"}}})
	})

}
