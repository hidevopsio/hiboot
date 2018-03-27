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
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hi/boot/pkg/log"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestUserLogin(t *testing.T) {
	baseUrl :=  os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")

	u := new(User)
	token, message, err := u.Login(baseUrl, username,  password)
	assert.Equal(t, nil, err)

	log.Debug(token)
	log.Debug(message)
}

func TestUserLoginFailed(t *testing.T) {
	baseUrl :=  os.Getenv("SCM_URL")

	u := new(User)
	token, message, err := u.Login(baseUrl, "xxx",  "xxx")
	assert.Contains(t, err.Error(), "Unauthorized")
	log.Debug(token)
	log.Debug(message)
}

func TestUserGetSession(t *testing.T) {
	baseUrl :=  os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")

	u := new(User)
	err := u.GetSession(baseUrl, username,  password)
	assert.Equal(t, nil, err)

	log.Debug(u.session)
}

func TestUserGetSessionUnauthorized(t *testing.T) {
	baseUrl :=  os.Getenv("SCM_URL")
	u := new(User)
	err := u.GetSession(baseUrl, "xxx",  "xxx")
	assert.Contains(t, err.Error(), "Unauthorized")

	log.Debug(err.Error())
}
