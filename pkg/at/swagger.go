package at

// Swagger is the annotation group - swagger
type Swagger struct {
	Annotation

	BaseAnnotation
}

// Operation Describes an operation or typically a HTTP method against a specific path.
// Operations with equivalent paths are grouped in a single Operation Object. A combination of a HTTP method and a path
// creates a unique operation.
// example:
// func (c *) CreateEmployee(at struct{
//     at.PostMapping `value:"/"`
//     at.Operation   `id:"Create Employee" description:"This is the employee creation api"`
//   }, request EmployeeRequest) {
//
//   ...
// }
type Operation struct {
	Annotation

	Swagger

	AtID          string `at:"id" json:"-"`
	AtDescription string `at:"description" json:"-"`
	AtSummary     string `at:"summary" json:"-"`
	AtDeprecated  bool   `at:"deprecated" json:"-"`
}

// Tags
type Tags struct {
	Annotation

	Swagger

	AtValues []string `at:"values" json:"-"`
}

// ExternalDocs
type ExternalDocs struct {
	Annotation

	Swagger

	AtDescription string `at:"description" json:"-"`
	AtURL         string `at:"url" json:"-"`
}

// ApiParam annotation to add additional meta-data for operation parameters
// func (c *) CreateEmployee(at struct{
//     at.PostMapping `value:"/"`
//     at.Parameter `value:"Employee object store in database table" required:"true"`
//   }, request EmployeeRequest) {
//
//   ...
// }
// Parameter
type Parameter struct {
	Annotation

	Swagger

	AtName        string `at:"name" json:"-"`
	AtType        string `at:"type" json:"-"`
	AtIn          string `at:"in" json:"-"`
	AtDescription string `at:"description" json:"-"`
}

// Produces corresponds to the `produces` field of the operation.
// Takes in comma-separated values of content types. For example, "application/json, application/xml" would suggest this
// operation generates JSON and XML output.
// example:
// at struct {
//    at.Consumes    `values:"application/json,application/xml"`
// }
type Produces struct {
	Annotation

	Swagger

	AtValues []string `at:"values" json:"-"`
}

// Consumes corresponds to the `consumes` field of the operation.
// Takes in comma-separated values of content types. For example, "application/json, application/xml" would suggest this
// API Resource accepts JSON and XML input.
// example:
// at struct {
//    at.Consumes    `values:"application/json,application/xml"`
// }
type Consumes struct {
	Annotation

	Swagger

	AtValues []string `at:"values" json:"-"`
}

// Response is the response type of the operation.
// example:
//Responses struct {
//	StatusOK struct {
//		at.Response `code:"200" description:"returns a greeting"`
//		at.Schema   `type:"string" description:"contains the actual greeting as plain text"`
//	}
//	StatusNotFound struct {
//		at.Response `code:"404" description:"greeter is not available"`
//		at.Schema   `type:"string" description:"Report 'not found' error message"`
//	}
//}
type Response struct {
	Annotation

	Swagger

	AtCode        int    `at:"code" json:"-"`
	AtDescription string `at:"description" json:"-"`
}

// Schema is the annotation that annotate Response or Parameter's properties
// example:
//Responses struct {
//	StatusOK struct {
//		at.Response `code:"200" description:"returns a greeting"`
//		at.Schema   `type:"string" description:"contains the actual greeting as plain text"`
//	}
//	StatusNotFound struct {
//		at.Response `code:"404" description:"greeter is not available"`
//		at.Schema   `type:"string" description:"Report 'not found' error message"`
//	}
//}
type Schema struct {
	Annotation

	Swagger

	AtType        string `at:"type" json:"-"`
	AtDescription string `at:"description" json:"-"`
}

// Header is the annotation that annotate the header
type Header struct {
	Annotation

	Swagger

	AtType        string `at:"type" json:"-"`
	AtFormat      string `at:"format" json:"-"`
	AtDescription string `at:"description" json:"-"`
}
