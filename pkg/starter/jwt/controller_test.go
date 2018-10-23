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
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/starter/jwt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
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
	token    jwt.Token
	tokenStr string
}

// for test only: token will expired in 1 second
var tokenExpiredSecond = int64(1)

func init() {
	log.SetLevel(log.DebugLevel)
}

func newFooController(token jwt.Token) *fooController {
	return &fooController{
		token: token,
	}
}

func (c *fooController) Get() string {
	log.Debug("fooController.Get")

	return "Hello, world"
}

func (c *fooController) PostLogin(request *userRequest) (response model.Response, err error) {
	log.Debug("fooController.Login")

	// you make validate username and password first
	token, err := c.token.Generate(jwt.Map{
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

func newBarController() *barController {
	return &barController{}
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

func (c *barController) Options() {
	log.Debug("barController.Options")
}

func TestJwtController(t *testing.T) {
	testApp := web.NewTestApplication(t, newFooController, newBarController)
	ctx := testApp.(app.ApplicationContext)

	fc := ctx.GetInstance(fooController{})
	assert.NotEqual(t, nil, fc)
	fooCtrl := fc.(*fooController)

	bc := ctx.GetInstance(barController{})
	assert.NotEqual(t, nil, bc)
	barCtrl := bc.(*barController)

	t.Run("should failed to get jwt properties when app is not run", func(t *testing.T) {
		_, ok := barCtrl.JwtProperties()
		assert.Equal(t, false, ok)

		_, ok = barCtrl.JwtPropertiesString()
		assert.Equal(t, false, ok)
	})

	t.Run("should login with POST /foo/login", func(t *testing.T) {
		testApp.
			Post("/foo/login").
			WithJSON(&userRequest{Username: "johndoe", Password: "iHop91#15"}).
			Expect().Status(http.StatusOK)
	})

	t.Run("should return http.StatusUnauthorized after GET /bar", func(t *testing.T) {
		testApp.Get("/bar").
			Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should return http.StatusOK after GET /bar", func(t *testing.T) {
		testApp.Options("/bar").
			Expect().Status(http.StatusOK)
	})

	token := fmt.Sprintf("Bearer %v", fooCtrl.tokenStr)
	t.Run("should return http.StatusOK on /bar with jwt token", func(t *testing.T) {

		testApp.Get("/bar").
			WithHeader("Authorization", token).
			Expect().Status(http.StatusOK)

		_, ok := barCtrl.JwtProperties()
		assert.Equal(t, true, ok)

		_, ok = barCtrl.JwtPropertiesString()
		assert.Equal(t, true, ok)

		username := barCtrl.JwtProperty("username")
		assert.Equal(t, "johndoe", username)

		testApp.Get("/bar").
			WithHeader("Authorization", "1df%^!"+token).
			Expect().Status(http.StatusUnauthorized)

		testApp.Get("/bar").
			WithHeader("Authorization", "").
			Expect().Status(http.StatusUnauthorized)
	})

	applicationContext := testApp.(app.ApplicationContext)
	jwtMiddleware := applicationContext.GetInstance(jwt.Middleware{}).(*jwt.Middleware)

	t.Run("should report error with SigningMethod", func(t *testing.T) {
		jwtMiddleware.Config.SigningMethod = jwtgo.SigningMethodRS512
		testApp.Get("/bar").
			WithHeader("Authorization", token).
			Expect().Status(http.StatusUnauthorized)
		jwtMiddleware.Config.SigningMethod = jwtgo.SigningMethodRS256
	})

	t.Run("should failed with invalid token", func(t *testing.T) {
		testApp.Get("/bar").
			WithHeader("Authorization", strings.Replace(token, "1", "2", -1)).
			Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should set CredentialsOptional", func(t *testing.T) {
		time.Sleep(2 * time.Second)

		testApp.Get("/bar").
			WithHeader("Authorization", token).
			Expect().Status(http.StatusUnauthorized)
	})

	t.Run("should set CredentialsOptional", func(t *testing.T) {
		jwtMiddleware.Config.CredentialsOptional = true
		testApp.Get("/bar").
			WithHeader("Authorization", "").
			Expect().Status(http.StatusOK)
	})

}

func TestAppWithoutJwtController(t *testing.T) {
	fooCtrl := new(fooController)
	testApp := web.NewTestApplication(t, fooCtrl)
	t.Run("should return http.StatusUnauthorized after GET /bar", func(t *testing.T) {
		testApp.Get("/foo").
			Expect().Status(http.StatusOK)
	})
}

func TestParseToken(t *testing.T) {
	claims := jwtgo.MapClaims{"username": "john"}
	jc := &jwt.Controller{}
	t.Run("should get username from jwt token", func(t *testing.T) {
		username := jc.ParseToken(claims, "username")
		assert.Equal(t, "john", username)
	})

	t.Run("should get empty string from jwt token", func(t *testing.T) {
		nonExist := jc.ParseToken(claims, "non-exist")
		assert.Equal(t, "", nonExist)
	})

}
