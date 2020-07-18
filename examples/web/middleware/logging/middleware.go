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
	ann := annotation.GetAnnotation(ctx.Annotations(), at.RequiresPermissions{})
	if ann != nil {
		va := ann.Field.Value.Interface().(at.RequiresPermissions)

		if va.AtType == "pagination" {
			log.Debugf("page number: %v, page size: %v", ctx.URLParam(va.AtIn[0]), ctx.URLParam(va.AtIn[1]))

			ctx.SetURLParam(va.AtOut[0], "where in(1,3,5,7)")
			ctx.SetURLParam(va.AtOut[0], "where in(2,4,6,8)") // just for test
		}

		log.Infof("[ logging middleware] %v - %v", ctx.GetCurrentRoute(), va.AtValues)
	}



	ann = annotation.GetAnnotation(ctx.Annotations(), at.Operation{})
	if ann != nil {
		va := ann.Field.Value.Interface().(at.Operation)
		log.Infof("[ logging middleware] %v - %v", ctx.GetCurrentRoute(), va.AtDescription)
	} else {
		log.Infof("[ logging middleware] %v", ctx.GetCurrentRoute())
	}

	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}

// PostLogging is the middleware post handler
func (m *loggingMiddleware) PostLogging( a struct{at.MiddlewarePostHandler `value:"/" `}, ctx context.Context) {
	ann := annotation.GetAnnotation(ctx.Annotations(), at.Operation{})
	if ann != nil {
		va := ann.Field.Value.Interface().(at.Operation)
		log.Infof("[post logging middleware] %v - %v - %v", ctx.GetCurrentRoute(), ctx.GetCurrentRoute(), va.AtDescription)
	} else {
		log.Infof("[post logging middleware] %v", ctx.GetCurrentRoute())
	}

	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}

