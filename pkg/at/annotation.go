package at

// Annotation is an empty struct that indicates the struct as an annotation
type Annotation struct {
}

// BaseAnnotation is the base of an annotation
type BaseAnnotation struct {
	Annotation

	Value string `json:"value" atf:""`
}
