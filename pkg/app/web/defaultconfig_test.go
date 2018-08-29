package web

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfiguration(t *testing.T) {
	dc := DefaultConfiguration()
	assert.Equal(t, "UTF-8", dc.Charset)
}
