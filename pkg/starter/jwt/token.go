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
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"io/ioutil"
	"github.com/hidevopsio/hiboot/pkg/log"
)

// JwtMap is the JWT map
type Map map[string]interface{}

type Token interface {
	Generate(payload Map, expired int64, unit time.Duration) (string, error)
}

type jwtToken struct {
	verifyKey     *rsa.PublicKey
	signKey       *rsa.PrivateKey
	//jwtMiddleware *JwtMiddleware
	jwtEnabled    bool
}

func NewJwtToken(p *Properties) Token {
	jt := new(jwtToken)
	jt.Initialize(p)
	return jt
}

func (t *jwtToken) Initialize(p *Properties) error {
	if io.IsPathNotExist(p.PrivateKeyPath) {
		log.Fatalf("private key file %v does not exist", p.PrivateKeyPath)
	}

	signBytes, err := ioutil.ReadFile(p.PrivateKeyPath)
	if err != nil {
		return err
	}

	t.signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}

	verifyBytes, err := ioutil.ReadFile(p.PublicKeyPath)
	if err != nil {
		return err
	}

	t.verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

	if err != nil {
		return err
	}

	t.jwtEnabled = true

	return nil
}

// Generate generates JWT token with specified exired time
func (t *jwtToken) Generate(payload Map, expired int64, unit time.Duration) (string, error) {
	if t.jwtEnabled {
		claim := jwt.MapClaims{
			"exp": time.Now().Add(unit * time.Duration(expired)).Unix(),
			"iat": time.Now().Unix(),
		}

		for k, v := range payload {
			claim[k] = v
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(t.signKey)

		return tokenString, err
	}

	return "", fmt.Errorf("JWT does not work")
}
