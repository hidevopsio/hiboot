package web

import (
	"testing"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestParseToken(t *testing.T) {
	claims := jwt.MapClaims{"username": "john"}
	jc := &JwtController{}
	username := jc.ParseToken(claims, "username")
	assert.Equal(t, "john", username)
}
