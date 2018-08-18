package aes

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	originalText := "encrypt this golang text"
	fmt.Println(originalText)
	key := []byte("example key 1234")

	// encrypt value to base64
	var cryptoText string
	t.Run("should encrypt with aes protocol", func(t *testing.T) {
		cryptoText = Encrypt(key, originalText)
		assert.NotEqual(t, "", cryptoText)
	})

	// encrypt base64 crypto to original value
	t.Run("should decrypt with aes protocol", func(t *testing.T) {
		text := Decrypt(key, cryptoText)
		assert.Equal(t, originalText, text)
	})
}
