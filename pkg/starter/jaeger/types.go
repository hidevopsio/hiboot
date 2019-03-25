package jaeger

import (
	"github.com/opentracing/opentracing-go"
)

//Tracer is the wrap of opentracing.Tracer
type Tracer opentracing.Tracer
