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
package controllers


import (
	"testing"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"fmt"
	"time"
)

func init() {
	utils.ChangeWorkDir("../")
}

func TestBarWithToken(t *testing.T) {
	app := web.NewTestApplication(t, new(BarController))

	pt, err := web.GenerateJwtToken(web.JwtMap{
		"username": "johndoe",
		"password": "PA$$W0RD",
	}, 100, time.Millisecond)
	if err == nil {

		t := fmt.Sprintf("Bearer %v", string(*pt))

		app.Get("/bar").
			WithHeader("Authorization", t).
			Expect().Status(http.StatusOK)

		time.Sleep(1 * time.Second)

		app.Get("/bar").
			WithHeader("Authorization", t).
			Expect().Status(http.StatusUnauthorized)
	}
}
