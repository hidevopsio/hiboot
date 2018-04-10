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

package controllers

import (
	"github.com/kataras/iris"

	"github.com/dgrijalva/jwt-go"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"github.com/hidevopsio/hi/cicd/pkg/ci/factories"
	"github.com/hidevopsio/hi/cicd/pkg/ci"
	"fmt"
	"strings"
	"github.com/hidevopsio/hi/boot/pkg/model"
)

type CicdRequest struct {
	Project  string `json:"project"` // Project = Namespace
	App      string `json:"app"`
	Profile  string `json:"profile"`
	Pipeline string `json:"pipeline"`
}

type CicdResponse struct {
	model.Response
}

// Operations about object
type CicdController struct{}

func init() {
	log.SetLevel(log.DebugLevel)
}

const (
	ScmUrl      = "url"
	ScmUsername = "username"
	ScmPassword = "password"
)

func (c *CicdController) Before(ctx iris.Context) {
	ctx.Application().Logger().Infof("Path: %s | IP: %s", ctx.Path(), ctx.RemoteAddr())

	// .Next is required to move forward to the chain of handlers,
	// if missing then it stops the execution at this handler.
	ctx.Next()
}

// @Title Deploy
// @Description deploy application by the pipeline
// @Param	body
// @Success 200 {string}
// @Failure 403 body is empty
// @router / [post]
func (c *CicdController) Run(ctx iris.Context) {
	log.Debug("CicdController.Run()")
	var pl ci.Pipeline
	err := ctx.ReadJSON(&pl)
	if err != nil {
		ctx.Values().Set("error", "deployment failed, read and parse request body failed. "+err.Error())
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	// decrypt jwt token
	token := ctx.Values().Get("jwt").(*jwt.Token)
	var username, password string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		pl.Scm.Url = parseToken(claims, ScmUrl)
		username = parseToken(claims, ScmUsername)
		password = parseToken(claims, ScmPassword)

		log.Debugf("url: %v, username: %v, password: %v", pl.Scm.Url, username, strings.Repeat("*", len(password)))

	} else {
		log.Debug(err)
	}

	// verify scm token
	// TODO:

	// invoke models
	pipelineFactory := new(factories.PipelineFactory)
	pipeline, err := pipelineFactory.New(pl.Name)
	message := "Successful."
	if err == nil {
		// Run Pipeline, password is a token, no need to pass username to pipeline
		pipeline.Init(&pl)
		err = pipeline.Run(username, password, false)
		if err != nil {
			message = err.Error()
		}
	} else {
		message = "Failed, " + err.Error()
	}
	response := &CicdResponse{
		Response: model.Response{
			Message: message,
		},
	}

	// just for debug now
	ctx.JSON(response)
}

func parseToken(claims jwt.MapClaims, prop string) string {
	return fmt.Sprintf("%v", claims[prop])
}
