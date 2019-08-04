package at

// StringAnnotation is the string annotation
type StringAnnotation interface {
	Value(value string) StringAnnotation
	String() (value string)
}

// IntAnnotation is the string annotation
type IntAnnotation interface {
	Int() (value int)
	Value(value int) IntAnnotation
}