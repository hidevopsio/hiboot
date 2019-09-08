package web

// HttpMethodSubscriber
type HttpMethodSubscriber interface {
	Subscribe(atController *Annotations, atMethod *Annotations)
}