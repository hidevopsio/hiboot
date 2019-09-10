package web

import "hidevops.io/hiboot/pkg/inject/annotation"

// HttpMethodSubscriber
type HttpMethodSubscriber interface {
	Subscribe(atController *annotation.Annotations, atMethod *annotation.Annotations)
}