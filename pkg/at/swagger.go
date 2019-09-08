package at

import "github.com/go-openapi/spec"

// Swagger annotation to declare swagger config
type Swagger struct {
	Annotation
}

// ApiOperation
// func (c *) CreateEmployee(at struct{
//     at.PostMapping `value:"/"`
//     at.ApiOperation `value:"createEmployee"`
//   }, request EmployeeRequest) {
//
//   ...
// }
type ApiOperation struct {
	// Required Element Summary
	Swagger

	Key string `value:"swagger.paths.${http.path}.${http.method}"`
	// Optional Element Summary
	spec.OperationProps `mapstructure:",squash"`
}

// Api is the annotation for REST Endpoints
// e.g.
// type employeeController struct {
//   at.RestController
//   at.Api `value:"Employee Management System" description:"Operations pertaining to employee in Employee Management System"`
// }
type Api struct {
	// Required Element Summary
	Swagger
}

// ApiParam annotation to add additional meta-data for operation parameters
// func (c *) CreateEmployee(at struct{
//     at.PostMapping `value:"/"`
//     at.ApiParam `value:"Employee object store in database table" required:"true"`
//   }, request EmployeeRequest) {
//
//   ...
// }
type ApiParam struct {
	// Optional Element Summary
	Swagger

	Required bool
}

// ApiResponse annotation to document other responses, in addition to the regular HTTP 200 OK, like this.
// func (c *) CreateEmployee(at struct{
//     at.PostMapping  `value:"/"`
//     at.ApiOperation `value:"Add an employee"`
//     at.ApiResponse  `200:"Successfully retrieved list" 401:"You are not authorized to view the resource 403:"Accessing the resource you were trying to reach is forbidden" 404:"The resource you were trying to reach is not found"`
//   }, request EmployeeRequest) (response Response) {
//
//   ...
// }
type ApiResponse struct {
	Swagger

	Key string `json:"key" value:"responses"`
	Code int `json:"code"`
	spec.ResponseProps
}

type ApiResponseSchema struct {
	Swagger

	Key string `value:"responses.${code}.schema"`
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