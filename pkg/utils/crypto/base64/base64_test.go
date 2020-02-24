package base64

import (
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/log"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestBase64String(t *testing.T) {
	src := "Hello world"
	str := EncodeToString(src)
	decoded, err := DecodeToString(str)
	assert.Equal(t, nil, err)
	assert.Equal(t, src, decoded)
	log.Debugf("src: %v, base64: %v, decoded: %v", src, str, decoded)
	assert.Equal(t, src, decoded)
}

func TestBase64Bytes(t *testing.T) {
	src := []byte("Abc123456")
	str := Encode(src)
	decoded, err := Decode(str)
	assert.Equal(t, nil, err)
	assert.Equal(t, src, decoded)
	log.Debugf("src: %v, base64: %v, decoded: %v", string(src), string(str), string(decoded))
	assert.Equal(t, src, decoded)
}
