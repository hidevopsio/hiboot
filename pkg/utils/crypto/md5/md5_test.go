package md5

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMd5(t *testing.T) {
	src := "Hello world"
	str := Encrypt(src)
	assert.Equal(t, "3e25960a79dbc69b674cd4ec67a72c62", str)
}
