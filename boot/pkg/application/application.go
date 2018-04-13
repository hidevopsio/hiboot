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


package application

import (
	"github.com/kataras/iris"
	log "github.com/kataras/golog"
	"github.com/hidevopsio/hi/boot/pkg/system"
	"fmt"
	"github.com/kataras/iris/context"
	"github.com/hidevopsio/hi/boot/pkg/utils"
	"crypto/rsa"
	"io/ioutil"
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"time"
	"github.com/hidevopsio/hi/boot/pkg/model"
)


const (
	privateKeyPath = "/config/app.rsa"
	pubKeyPath     = "/config/app.rsa.pub"
)

type Boot struct {
	app    *iris.Application
	config *system.Configuration
}

type MapJwt map[string]interface{}

type JwtToken string

type Health struct {
	Status string `json:"status"`
}

var (
	jwtHandler *jwtmiddleware.Middleware
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
	sysCfg *system.Configuration
)


const (
	application = "application"
	config = "/config"
	yaml = "yaml"
)


func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {

	b := &system.Builder{
		Path: utils.GetWorkingDir("boot/pkg/system/builder_test.go") + config,
		Name: application,
		FileType: yaml,
		Profile: "local",
		ConfigType: system.Configuration{},
	}
	sysCfg := cp.(*system.Configuration)
	log.Print(sysCfg)

	log.SetLevel(sysCfg.Logging.Level)

	wd := utils.GetWorkingDir("/boot/pkg/application/application.go")

	signBytes, err := ioutil.ReadFile(wd + privateKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(wd + pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

func GenerateJwtToken(payload MapJwt, expired int64, unit time.Duration) (JwtToken, error)  {

	claim := jwt.MapClaims{
		"exp": time.Now().Add(unit * time.Duration(expired)).Unix(),
		"iat": time.Now().Unix(),
	}

	for k, v := range payload {
		claim[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(signKey)

	jwtToken := JwtToken(tokenString)

	return jwtToken, err
}

// set response
func Response(ctx iris.Context, message string, data interface{})  {
	response := &model.Response{
		Code: ctx.GetStatusCode(),
		Message: message,
		Data: data,
	}

	// just for debug now
	ctx.JSON(response)
}

// set response
func ResponseError(ctx iris.Context, message string, code int)  {
	response := &model.Response{
		Code: code,
		Message: message,
	}

	// just for debug now
	ctx.StatusCode(code)
	ctx.JSON(response)
}


func (b *Boot) Init() {
	jwtHandler = jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			//log.Debug(token)
			return verifyKey, nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodRS256,
	})

	b.config = sysCfg
	log.Debug(sysCfg)
	log.Debug("application init")

	b.app = iris.New()

	b.app.Get("/health", func(ctx context.Context) {
		health := Health{
			Status: "UP",
		}
		ctx.JSON(health)
	})
}


func (b *Boot) App() *iris.Application {
	return b.app
}

func (b *Boot) Config() *system.Configuration  {
	return b.config
}

func (b *Boot) GetSignKey() *rsa.PrivateKey  {
	return signKey
}

func (b *Boot) ApplyJwt()  {
	b.app.Use(jwtHandler.Serve)
}

func (b *Boot) Run() {
	serverPort := fmt.Sprintf(":%v", b.config.Server.Port)
	b.app.Run(iris.Addr(fmt.Sprintf(serverPort)), iris.WithCharset("UTF-8"), iris.WithoutVersionChecker)
}