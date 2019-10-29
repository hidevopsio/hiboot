package at

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestControllerAnnotations(t *testing.T) {

	t.Run("should pass test for annotation RequestMapping", func(t *testing.T) {
		rm := RequestMapping{HttpMethod: HttpMethod{BaseAnnotation: BaseAnnotation{AtValue: "/foo"}}}

		assert.Equal(t, rm.BaseAnnotation.AtValue, "/foo")
		assert.Equal(t, rm, RequestMapping{HttpMethod: HttpMethod{BaseAnnotation: BaseAnnotation{AtValue: "/foo"}}})
	})

}
