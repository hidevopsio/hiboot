package at

import "github.com/go-openapi/spec"

// Swagger annotation to declare swagger config
type Swagger struct {
	Annotation

	BaseAnnotation
}

// Operation
// func (c *) CreateEmployee(at struct{
//     at.PostMapping `value:"/"`
//     at.Operation `value:"createEmployee"`
//   }, request EmployeeRequest) {
//
//   ...
// }
type Operation struct {
	Annotation

	// Required Element Summary
	Swagger

	// Optional Element Summary
	spec.Operation
}

// OpenAPIDefinition is the annotation for swagger
type OpenAPIDefinition struct {
	Annotation

	Swagger
}

// Schemes is the annotation for Swagger
type Schemes struct {
	Annotation

	Swagger

	Values []string `json:"values"`
}


// ApiParam annotation to add additional meta-data for operation parameters
// func (c *) CreateEmployee(at struct{
//     at.PostMapping `value:"/"`
//     at.Parameter `value:"Employee object store in database table" required:"true"`
//   }, request EmployeeRequest) {
//
//   ...
// }

// ParameterItem
type ParameterItem struct {
	Annotation

	Swagger
}

// Parameter
type Parameter struct {
	Annotation

	ParameterItem

	spec.Parameter
}

// Produces
type Produces struct{
	Annotation

	Swagger

	Values []string `json:"values"`
}

// Consumes
type Consumes struct{
	Annotation

	Swagger

	Values []string `json:"values"`
}

// Response annotation to document other responses, in addition to the regular HTTP 200 OK, like this.
// func (c *) CreateEmployee(at struct{
//     at.PostMapping  `value:"/"`
//     at.Operation `value:"Add an employee"`
//     at.Response  `200:"Successfully retrieved list" 401:"You are not authorized to view the resource 403:"Accessing the resource you were trying to reach is forbidden" 404:"The resource you were trying to reach is not found"`
//   }, request EmployeeRequest) (response Response) {
//
//   ...
// }
type Response struct {
	Annotation

	Swagger

	Code int `json:"code"`
	spec.Response
}

type Schema struct {
	Annotation

	Swagger

	spec.Schema
}
