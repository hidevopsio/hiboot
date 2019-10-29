package webutils

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/inject/annotation"
	"hidevops.io/hiboot/pkg/log"
	"testing"
)

func TestUtils(t *testing.T) {
	testData := []interface{}{
		&struct {
			at.GetMapping `value:"/"`
		}{},
		&struct {
			at.PostMapping `value:"/"`
		}{},
		&struct {
			at.PutMapping `value:"/"`
		}{},
		&struct {
			at.DeleteMapping `value:"/"`
		}{},
		&struct {
			at.PatchMapping `value:"/"`
		}{},
		&struct {
			at.OptionsMapping `value:"/"`
		}{},
		&struct {
			at.AnyMapping `value:"/"`
		}{},
		&struct {
			at.TraceMapping `value:"/"`
		}{},
	}

	for _, a := range testData {
		ann := annotation.GetAnnotations(a)
		err := annotation.Inject(ann.Items[0])
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(ann.Items))
		method, path := GetHttpMethod(ann)
		log.Debug(method)
		log.Debug(path)
		assert.Equal(t, "/", path)
	}
}
