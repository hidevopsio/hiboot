package logging

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/inject/annotation"
	"hidevops.io/hiboot/pkg/log"
)

type loggingMiddleware struct {
	at.Middleware
}

func newLoggingMiddleware() *loggingMiddleware {
	return &loggingMiddleware{}
}

func init() {
	app.Register(newLoggingMiddleware)
}

// Logging is the middleware handler,it support dependency injection, method annotation
// middleware handler can be annotated to specific purpose or general purpose
func (m *loggingMiddleware) Logging( a struct{at.MiddlewareHandler `value:"/" `}, ctx context.Context) {
	atRequiresPermissions := annotation.GetAnnotation(ctx.Annotations(), at.RequiresPermissions{})
	if atRequiresPermissions != nil {
		log.Infof("[ logging middleware] %v - %v", ctx.GetCurrentRoute(), atRequiresPermissions.Field.StructField.Tag.Get("value"))
	}

	atOperation := annotation.GetAnnotation(ctx.Annotations(), at.Operation{})
	if atOperation != nil {
		log.Infof("[ logging middleware] %v - %v", ctx.GetCurrentRoute(), atOperation.Field.StructField.Tag.Get("description"))
	} else {
		log.Infof("[ logging middleware] %v", ctx.GetCurrentRoute())
	}

	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}

// PostLogging is the middleware post handler
func (m *loggingMiddleware) PostLogging( a struct{at.MiddlewarePostHandler `value:"/" `}, ctx context.Context) {
	atOperation := annotation.GetAnnotation(ctx.Annotations(), at.Operation{})
	if atOperation != nil {
		log.Infof("[post logging middleware] %v - %v", ctx.GetCurrentRoute(), atOperation.Field.StructField.Tag.Get("description"))
	} else {
		log.Infof("[post logging middleware] %v", ctx.GetCurrentRoute())
	}

	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}

