// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jwt

import (
	"time"
	jwtgo "github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"crypto/rsa"
	"io/ioutil"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/system"
	"fmt"
)

type Map map[string]interface{}

type Token string

const (
	privateKeyPath = "/config/app.rsa"
	pubKeyPath     = "/config/app.rsa.pub"
)

var (
	verifyKey  *rsa.PublicKey
	signKey    *rsa.PrivateKey
	jwtHandler *jwtmiddleware.Middleware
	jwtEnabled bool
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Init(wd string) error  {

	// check if key exist
	if utils.IsPathNotExist(wd + privateKeyPath) {
		return &system.NotFoundError{Name: wd + privateKeyPath}
	}


	signBytes, err := ioutil.ReadFile(wd + privateKeyPath)
	if err != nil {
		return err
	}

	signKey, err = jwtgo.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}

	verifyBytes, err := ioutil.ReadFile(wd + pubKeyPath)
	if err != nil {
		return err
	}

	verifyKey, err = jwtgo.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return err
	}

	jwtHandler = jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwtgo.Token) (interface{}, error) {
			//log.Debug(token)
			return verifyKey, nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwtgo.SigningMethodRS256,
	})

	jwtEnabled = true

	return nil
}


func GetSignKey() *rsa.PrivateKey {
	return signKey
}

func GetHandler() *jwtmiddleware.Middleware  {
	return jwtHandler
}

func GenerateToken(payload Map, expired int64, unit time.Duration) (*Token, error) {
	if jwtEnabled {
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

	return nil, fmt.Errorf("JWT does not work")
}
