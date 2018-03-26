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

	"github.com/hidevopsio/hi/cicd/pkg/pipeline"
	"github.com/hidevopsio/hi/cicd/pkg/ci"
)

// Operations about object
type CicdController struct {

}

func (c *CicdController) Before(ctx iris.Context) {
	ctx.Application().Logger().Infof("Path: %s | IP: %s", ctx.Path(), ctx.RemoteAddr())

	// .Next is required to move forward to the chain of handlers,
	// if missing then it stops the execution at this handler.
	ctx.Next()
}


// @Title Deploy
// @Description deploy application
// @Param	body
// @Success 200 {string}
// @Failure 403 body is empty
// @router / [post]
func (c *CicdController) Run(ctx iris.Context) {
	var pipeline pipeline.Pipeline
	err := ctx.ReadJSON(&pipeline)
	if err != nil {
		ctx.Values().Set("error", "deployment failed, read and parse request body failed. " + err.Error())
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	// invoke models
	ci.Run(&pipeline)

	// just for debug now
	ctx.JSON(pipeline)
}
