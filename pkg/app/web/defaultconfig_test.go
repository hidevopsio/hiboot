package web

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultConfiguration(t *testing.T) {
	dc := DefaultConfiguration()
	assert.Equal(t, "UTF-8", dc.Charset)
}
