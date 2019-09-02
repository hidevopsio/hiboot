package at

// Middleware is the annotation that annotate the controller or method use middleware
type Middleware struct {
	Annotation
}

// MiddlewareHandler is the annotation that annotate the controller or method use middleware
type MiddlewareHandler struct {
	Annotation
}

// UseMiddleware is the annotation that that annotate the controller or method use middleware based on condition
type UseMiddleware struct {
	Conditional
}

// UseJwt
type UseJwt struct {
	UseMiddleware
}