package at

// Conditional check if the string value of any give condition of the struct
//
//	type Example struct {
//	  at.Conditional `value:"your-condition-express"`
//
//	  Name string
//	}
type Conditional struct {
	Annotation

	BaseAnnotation
}

// ConditionalOnField annotation check if the string value of give fields of the struct
//
//	type Example struct {
//	  at.ConditionalOnField `value:"Namespace,Name"`
//
//	  Namespace string
//	  Name string
//	}
type ConditionalOnField struct {
	Annotation

	BaseAnnotation
}
