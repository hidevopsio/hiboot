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

// Line 1: main package
package main

import (
	"hidevops.io/hiboot/pkg/app/web"
	"net/http"
	"testing"
	"time"
)

func TestRunMain(t *testing.T) {
	go main()
}

func TestController(t *testing.T) {
	time.Sleep(time.Second)
	testApp := web.NewTestApp(t, new(UserController), new(orgController)).SetProperty("server.context_path", "/router-example").Run(t)

	testApp.Post("/router-example/user").
		WithJSON(&UserRequest{UserVO: UserVO{Username: "John", Password: "password"}}).
		Expect().Status(http.StatusOK)

	testApp.Get("/router-example/user/123/and/John").
		Expect().Status(http.StatusOK)

	testApp.Delete("/router-example/user/123").
		Expect().Status(http.StatusOK)

	testApp.Get("/router-example/user").
		Expect().Status(http.StatusOK)

	testApp.Get("/router-example/organization/official-site").
		Expect().Status(http.StatusOK)

	testApp.Get("/router-example/organization/123/and/John").
		Expect().Status(http.StatusOK)
}
