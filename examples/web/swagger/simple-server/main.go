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
	"hidevops.io/hiboot/pkg/starter/swagger"
	"net/http"
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
	at.GetMapping     `value:"/{id:int}"`
	at.Operation      `value:"Get an employee"`
	at.Parameter      `value:"Path variable employee ID" required:"true"`
}, id int) (response *EmployeeResponse, err error) {
	response = new(EmployeeResponse)
	for _, e := range c.dummyData {
		if id == e.Id {
			response.Code = http.StatusOK
			response.Message = "success"
			response.Data = e
			break
		}
	}

	if response.Data == nil {
		response.Code = http.StatusNotFound
		response.Message = "Resource is not found"
		err = fmt.Errorf("employee %v is not found", id)
	}

	return
}

// ListEmployee
func (c *employeeController) ListEmployee(at struct{
	at.GetMapping     `value:"/"`
	at.Operation      `value:"List employees"`
	at.Parameter      `value:"Path variable employee ID" required:"true"`
	at.Response 	   `value:"Successfully list employee"`
}) (response *ListEmployeeResponse, err error) {
	response = new(ListEmployeeResponse)
	response.Code = http.StatusOK
	response.Message = "success"
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
			swagger.Profile,
		).Run()
}
