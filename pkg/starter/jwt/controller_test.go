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

package jwt_test

import (
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/starter/jwt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

type userRequest struct {
	model.RequestBody
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type fooController struct {
	web.Controller
	jwtToken jwt.Token
	tokenStr string
}

// for test only: token will expired in 1 second
var tokenExpiredSecond = int64(1)

func init() {
	log.SetLevel(log.DebugLevel)
}

func (c *fooController) Init(jwtToken jwt.Token) {
	c.jwtToken = jwtToken
}

func (c *fooController) Get() string {
	log.Debug("fooController.Get")

	return "Hello, world"
}

func (c *fooController) PostLogin(request *userRequest) (response model.Response, err error) {
	log.Debug("fooController.Login")

	// you make validate username and password first
	token, err := c.jwtToken.Generate(jwt.Map{
		"username": request.Username,
		"password": request.Password,
	}, tokenExpiredSecond, time.Second)
	response = new(model.BaseResponse)
	response.SetData(token)
	c.tokenStr = token

	return
}

// BarController
type barController struct {
	jwt.Controller
}

func (c *barController) Before() {
	log.Debug("barController.Before")

	_, ok := c.JwtProperties()
	// intercept all requests that not contain jwt token
	if !ok {
		return
	}
	c.Ctx.Next()
}

func (c *barController) Get() string {
	log.Debug("barController.Get")

	return "Hello, world"
}

func TestJwtController(t *testing.T) {
	fooCtrl := new(fooController)
	barCtrl := new(barController)

	t.Run("should failed to get jwt properties when app is not run", func(t *testing.T) {
		_, ok := barCtrl.JwtProperties()
		assert.Equal(t, false, ok)

		_, ok = barCtrl.JwtPropertiesString()
		assert.Equal(t, false, ok)
	})

	app := web.NewTestApplication(t, fooCtrl, barCtrl)

	t.Run("should login with POST /foo/login", func(t *testing.T) {
		app.
			Post("/foo/login").
			WithJSON(&userRequest{Username: "johndoe", Password: "iHop91#15"}).
			Expect().Status(http.StatusOK)
	})

	t.Run("should return http.StatusUnauthorized after GET /bar", func(t *testing.T) {
		app.Get("/bar").
			Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should return http.StatusOK on /bar with jwt token", func(t *testing.T) {
		token := fmt.Sprintf("Bearer %v", fooCtrl.tokenStr)

		app.Get("/bar").
			WithHeader("Authorization", token).
			Expect().Status(http.StatusOK)

		_, ok := barCtrl.JwtProperties()
		assert.Equal(t, true, ok)

		_, ok = barCtrl.JwtPropertiesString()
		assert.Equal(t, true, ok)

		username := barCtrl.JwtProperty("username")
		assert.Equal(t, "johndoe", username)

		time.Sleep(2 * time.Second)

		app.Get("/bar").
			WithHeader("Authorization", token).
			Expect().Status(http.StatusUnauthorized)
	})

}

func TestAppWithoutJwtController(t *testing.T) {
	fooCtrl := new(fooController)
	app := web.NewTestApplication(t, fooCtrl)
	t.Run("should return http.StatusUnauthorized after GET /bar", func(t *testing.T) {
		app.Get("/foo").
			Expect().Status(http.StatusOK)
	})
}

func TestParseToken(t *testing.T) {
	claims := jwtgo.MapClaims{"username": "john"}
	jc := &jwt.Controller{}
	username := jc.ParseToken(claims, "username")
	assert.Equal(t, "john", username)
}
