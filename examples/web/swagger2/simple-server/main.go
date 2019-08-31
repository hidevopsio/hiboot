package main

import (
	"fmt"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
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
	Data            *Employee `json:"data,omitempty" api:"The employee data"`
}

type ListEmployeeResponse struct {
	at.ResponseBody `json:"-"`
	at.ApiModel     `description:"All details about the Employee. " json:"-"`
	model.BaseResponse
	Data            []*Employee `json:"data,omitempty" api:"The employee data"`
}

type employeeController struct {
	at.RestController
	at.RequestMapping `value:"/employee"`
	at.Api `value:"Employee Management System" description:"Operations pertaining to employee in Employee Management System"`

	dummyData []*Employee
}

// newEmployeeController is the constructor for orgController
// you may inject dependency to constructor
func newEmployeeController() *employeeController {
	return &employeeController{
		dummyData: []*Employee{
			{
				Id:        123,
				FirstName: "John",
				LastName:  "Deng",
			},
			{
				Id:        456,
				FirstName: "Mike",
				LastName:  "Philip",
			},
		},
	}
}

// Before
func (c *employeeController) BeforeMethod(at struct{at.BeforeMethod}, ctx context.Context)  {
	log.Debug("before method")
	ctx.Next()
	return
}

// GetEmployee
// at.GetMapping is an annotation to define request mapping for http method GET,
func (c *employeeController) GetEmployee(at struct {
	at.PathVariable
	at.GetMapping `value:"/{id:int}"`
	at.ApiOperation `value:"Get an employee"`
	at.ApiParam `value:"Path variable employee ID" required:"true"`
	at.ApiResponse200 `value:"Successfully get an employee"`
	at.ApiResponse404 `value:"The resource you were trying to reach is not found"`
}, id int) (response *EmployeeResponse, err error) {
	response = new(EmployeeResponse)
	for _, e := range c.dummyData {
		if id == e.Id {
			response.Code = at.ApiResponse200.Code
			response.Message = at.ApiResponse200.Value
			response.Data = e
			break
		}
	}

	if response.Data == nil {
		response.Code = at.ApiResponse404.Code
		response.Message = at.ApiResponse404.Value
		err = fmt.Errorf("employee %v is not found", id)
	}

	return
}

// ListEmployee
func (c *employeeController) ListEmployee(at struct{
	at.GetMapping `value:"/"`
	at.ApiOperation `value:"List employees"`
	at.ApiParam `value:"Path variable employee ID" required:"true"`
	at.ApiResponse200 `value:"Successfully list employee"`
	at.ApiResponse404 `value:"The resource you were trying to reach is not found"`
}) (response *ListEmployeeResponse, err error) {
	response = new(ListEmployeeResponse)
	response.Code = at.ApiResponse200.Code
	response.Message = at.ApiResponse200.Value
	response.Data = c.dummyData
	return
}

// After
func (c *employeeController) AfterMethod(at struct{at.AfterMethod}, ctx context.Context)  {
	log.Debug("before method")
	ctx.Next()
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
