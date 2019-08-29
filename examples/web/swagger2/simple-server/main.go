package main

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/factory"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hiboot/pkg/starter/actuator"
	"hidevops.io/hiboot/pkg/starter/swagger2"
)


type Employee struct {

	Id        int    `api:"The database generated employee ID" json:"id"`
	FirstName string `api:"The employee first name" json:"first_name"`
	LastName  string `api:"The employee last name" json:"last_name"`
}

type EmployeeResponse struct {
	at.ResponseBody `json:"-"`
	at.ApiModel     `description:"All details about the Employee. " json:"-"`
	model.BaseResponse
	Data            Employee `json:"data" api:"The employee data"`
}


type employeeController struct {
	at.RestController
	at.RequestMapping `value:"/employee"`
	at.Api `value:"Employee Management System" description:"Operations pertaining to employee in Employee Management System"`
}

// newEmployeeController is the constructor for orgController
// you may inject dependency to constructor
func newEmployeeController() *employeeController {
	return new(employeeController)
}

// Get
// at.GetMapping is an annotation to define request mapping for http method GET,
func (c *employeeController) Get(at struct {
	at.PathVariable
	at.GetMapping `value:"/{id:int}"`
	at.ApiOperation `value:"Get an employee"`
	at.ApiParam `value:"Path variable employee ID" required:"true"`
	at.ApiResponse200 `value:"Successfully get an employee"`
	at.ApiResponse404 `value:"The resource you were trying to reach is not found"`
}, id int) (response *EmployeeResponse, err error) {
	log.Infof("annotations: %v", at)

	response = new(EmployeeResponse)
	response.Code = at.ApiResponse200.Code
	response.Message = at.ApiResponse200.Value
	response.Data = Employee{
		Id: id,
		FirstName: "John",
		LastName: "Deng",
	}
	return
}

// ListEmployee
func (c *employeeController) ListEmployee(at struct{ at.GetMapping `value:"/"` }, factory factory.ConfigurableFactory) (response model.Response, err error) {
	response = new(model.BaseResponse)
	return
}

func init() {
	app.Register(newEmployeeController)
}

// Hiboot main function
func main() {
	// create new web application and run it
	web.NewApplication().
		SetProperty(app.ProfilesInclude,
			actuator.Profile,
			swagger2.Profile,
		).Run()
}
