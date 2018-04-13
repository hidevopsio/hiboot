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


package auth

import (
	"github.com/hidevopsio/hi/cicd/pkg/scm/factories"
	"github.com/hidevopsio/hi/cicd/pkg/scm"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"crypto/rsa"
)

type UserInterface interface{
	Login(baseUrl, username, password string) (string, error)
}

type Credential struct{
	username string
	password string
	isToken bool
}

type User struct {
	session scm.SessionInterface
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (u *User) Login(baseUrl, username, password string) (string, string, error) {
	log.Debug("User.Login()")
	retVal := "Login successful."
	// login scm
	err := u.GetSession(baseUrl, username, password)
	if err != nil {
		retVal = "Login " + err.Error()
	}

	return u.session.GetToken(), retVal, err
}

func (u *User) GetSession(baseUrl, username, password string) error {
	log.Debug("User.GetSession()")
	scmFactory := new(factories.ScmFactory)
	var err error
	u.session, err = scmFactory.New(factories.GitlabScmType)

	if err == nil {
		err = u.session.GetSession(baseUrl, username, password)
	}

	return err
}
