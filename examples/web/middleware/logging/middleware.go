package logging

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
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
func (m *loggingMiddleware) Logging(at struct {
	at.MiddlewareHandler `value:"/" `
}, ctx context.Context) {

	log.Infof("[logging middleware] %v", ctx.GetCurrentRoute())

	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}
