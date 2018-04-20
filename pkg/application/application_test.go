package application

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/dgrijalva/jwt-go"
	"github.com/hidevopsio/hiboot/pkg/log"
	"fmt"
	"time"
	"io/ioutil"
	"github.com/hidevopsio/hiboot/pkg/utils"
)

func init() {
	log.SetLevel(log.DebugLevel)

	wd := utils.GetWorkingDir("")

	signBytes, err := ioutil.ReadFile(wd + privateKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(wd + pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

func TestGenerateJwtToken(t *testing.T) {
	jwtToken, err := GenerateJwtToken(MapJwt{"foo": "bar"}, 1, time.Hour)

	log.Debug(jwtToken)
	token, err := jwt.Parse(string(jwtToken), func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return verifyKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Debugf("%v, %v, %v", claims["foo"], claims["exp"], claims["iat"])
	} else {
		log.Debug(err)
	}

	assert.Equal(t, nil, err)
}
