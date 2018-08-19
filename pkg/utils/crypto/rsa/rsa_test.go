package rsa

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/utils/crypto"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestRsa(t *testing.T) {
	src := []byte("hello")
	data, _ := Encrypt([]byte(src))
	decrypted, err := Decrypt(data)
	assert.Equal(t, nil, err)
	log.Debugf("encrypted: %v, decrypted: %v, org: %v", string(data), string(decrypted), string(src))
	assert.Equal(t, src, decrypted)
}


func TestRsaBase64(t *testing.T) {
	src := []byte("hello")
	data, _ := EncryptBase64([]byte(src))
	decrypted, err := DecryptBase64(data)
	assert.Equal(t, nil, err)
	log.Debugf("encrypted: %v, decrypted: %v, org: %v", string(data), string(decrypted), string(src))
	assert.Equal(t, src, decrypted)
}

func TestExeptions(t *testing.T) {
	t.Run("should report error with invalid public key", func(t *testing.T) {
		src := []byte("hello")
		_, err := Encrypt([]byte(src), []byte("invalid-key"))
		assert.Equal(t, crypto.InvalidPublicKeyError, err)
	})

	t.Run("should report error with invalid private key", func(t *testing.T) {
		src := []byte("hello")
		data, _ := Encrypt([]byte(src))
		_, err := Decrypt(data, []byte("invalid-key"))
		assert.Equal(t, crypto.InvalidPrivateKeyError, err)
	})
}

