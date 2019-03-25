package jaeger

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	webctx "hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
)

//Span  is the wrap of opentracing.Span
type Span struct {
	at.ContextAware
	opentracing.Span
	context webctx.Context
}

//ChildSpan
type ChildSpan Span


func (s *Span) Inject(ctx context.Context, method string, url string, operationName string) opentracing.Span {
	c := opentracing.ContextWithSpan(ctx, s)
	newSpan, _ := opentracing.StartSpanFromContext(c, operationName)

	ext.SpanKindRPCClient.Set(newSpan)
	ext.HTTPUrl.Set(newSpan, url)
	ext.HTTPMethod.Set(newSpan, method)
	newSpan.Tracer().Inject(
		newSpan.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(s.context.Request().Header),
	)

	return newSpan
}