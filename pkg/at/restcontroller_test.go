package at

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestControllerAnnotations(t *testing.T) {

	t.Run("should pass test for annotation RequestMapping", func(t *testing.T) {
		rm := RequestMapping{Annotation{"/foo"}}

		assert.Equal(t, rm.Annotation.Value, "/foo")
		assert.Equal(t, rm, RequestMapping{Annotation{"/foo"}})
	})

	//t.Run("should pass test for annotation Method", func(t *testing.T) {
	//	rm := Method("foo")
	//
	//	assert.Equal(t, rm.String(), "foo")
	//	assert.Equal(t, rm, Method("foo"))
	//	assert.Equal(t, rm.Value("bar"), Method("bar"))
	//})
	//t.Run("should pass test for annotation Path", func(t *testing.T) {
	//	rm := Path("foo")
	//
	//	assert.Equal(t, rm.String(), "foo")
	//	assert.Equal(t, rm, Path("foo"))
	//	assert.Equal(t, rm.Value("bar"), Path("bar"))
	//})

}
