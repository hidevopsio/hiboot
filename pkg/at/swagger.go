package at

import "github.com/go-openapi/spec"

// Swagger annotation to declare swagger config
type Swagger struct {
	Annotation
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
	// Required Element Summary
	Swagger

	Key string `value:"swagger.paths.${at.http.path}.${at.http.method}"`
	// Optional Element Summary
	spec.Operation `mapstructure:",squash"`
}

// OpenAPIDefinition is the annotation for swagger
type OpenAPIDefinition struct {
	Swagger
}

// Info is the annotation for Swagger Info
type Info struct {
	Swagger

	Key string `json:"key" value:"swagger.info"`

	spec.Info
}


// Host is the annotation for Swagger
type Host struct {
	Swagger

	Key string `json:"key" value:"swagger.host"`
}

// BasePath is the annotation for Swagger
type BasePath struct {
	Swagger

	Key string `json:"key" value:"swagger.basePath"`
}

// Schemes is the annotation for Swagger
type Schemes struct {
	Swagger

	Key string `json:"key" value:"swagger.schemes"`

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
	Swagger
}

// Parameter
type Parameter struct {
	ParameterItem

	Key string `json:"key" value:"parameters"`
	
	spec.Parameter
}

// Produces
type Produces struct{
	Swagger

	Key string `json:"key" value:"produces"`

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
	Swagger

	Key string `json:"key" value:"responses.${at.response.code}"`
	Code int `json:"code"`
	spec.Response
}

type ResponseSchema struct {
	Swagger

	Key string `value:"responses.${at.response.code}.schema"`
	Code int `json:"code"`

	spec.Schema
}

// ApiModel annotation to describe the properties of the  Employee  model.
//type Employee struct {
//	ApiModel `description:"All details about the Employee. "`
//
//	Id int `api:"The database generated employee ID"`
//	FirstName string `api:"The employee first name"`
//	LastName string `api:"The employee last name"`
//}
type ApiModel struct {
	Swagger
	Description string
}