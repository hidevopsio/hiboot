package jwt

import (
	"time"
	jwtgo "github.com/dgrijalva/jwt-go"
)

type Map map[string]interface{}

type Token string

func GenerateToken(payload Map, expired int64, unit time.Duration, signKey interface{}) (*Token, error) {

	claim := jwtgo.MapClaims{
		"exp": time.Now().Add(unit * time.Duration(expired)).Unix(),
		"iat": time.Now().Unix(),
	}

	for k, v := range payload {
		claim[k] = v
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS256, claim)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(signKey)

	jwtToken := Token(tokenString)

	return &jwtToken, err
}
