package controllers

import (
	"github.com/kataras/iris"

	"github.com/hi-devops-io/hi/cicd/pkg/pipelines"
	"github.com/hi-devops-io/hi/cicd/pkg/config"
)

// Operations about object
type DeploymentController struct {

}

func (c *DeploymentController) Before(ctx iris.Context) {
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
func (c *DeploymentController) Deploy(ctx iris.Context) {
	var pipeline config.Pipeline
	err := ctx.ReadJSON(&pipeline)
	if err != nil {
		ctx.Values().Set("error", "deployment failed, read and parse request body failed. " + err.Error())
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	// invoke models
	pipelines.Deploy(&pipeline)

	// just for debug now
	ctx.JSON(pipeline)
}
