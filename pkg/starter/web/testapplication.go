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

package web

import (
	"testing"
	"net/http"
	"github.com/iris-contrib/httpexpect"
	"github.com/stretchr/testify/assert"
	"github.com/kataras/iris/httptest"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type TestApplicationInterface interface {
	ApplicationInterface
	NewTestServer(t *testing.T) *httpexpect.Expect
}

type TestApplication struct {
	Application
	expect *httpexpect.Expect
}

func NewTestApplication(t *testing.T, controllers... interface{}) *TestApplication {
	log.SetLevel(log.DebugLevel)
	ta := new(TestApplication)
	err := ta.Init(controllers...)
	assert.Equal(t, nil, err)
	ta.expect = ta.RunTestServer(t)
	return ta
}

func (wa *TestApplication) RunTestServer(t *testing.T) *httpexpect.Expect {
	return httptest.New(t, wa.app)
}

func (ta *TestApplication) Request(method, path string, pathargs ...interface{}) *httpexpect.Request {
	return ta.expect.Request(method, path, pathargs...)
}


func (ta *TestApplication) Post(path string, pathargs ...interface{}) *httpexpect.Request {
	return ta.expect.Request(http.MethodPost, path, pathargs...)
}

func (ta *TestApplication) Put(path string, pathargs ...interface{}) *httpexpect.Request {
	return ta.expect.Request(http.MethodPut, path, pathargs...)
}

func (ta *TestApplication) Patch(path string, pathargs ...interface{}) *httpexpect.Request {
	return ta.expect.Request(http.MethodPatch, path, pathargs...)
}

func (ta *TestApplication) Get(path string, pathargs ...interface{}) *httpexpect.Request {
	return ta.expect.Request(http.MethodGet, path, pathargs...)
}

func (ta *TestApplication) Delete(path string, pathargs ...interface{}) *httpexpect.Request {
	return ta.expect.Request(http.MethodDelete, path, pathargs...)
}