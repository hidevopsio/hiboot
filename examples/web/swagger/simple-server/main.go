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
	"hidevops.io/hiboot/pkg/starter/logging"
	"hidevops.io/hiboot/pkg/starter/swagger"
	"net/http"
)

type Assert struct {
	at.Schema    `json:"-"`
	ID int `schema:"The assert ID" json:"id"`
	Name string `schema:"The assert name" json:"name"`
}

type Manager struct {
	at.Schema `json:"-"`
	ID int `schema:"The manager ID" json:"id"`
	Name string `schema:"The manager name of the employee" json:"name"`
}

type Employee struct {
	at.Schema    `json:"-"`
	Id        int    `schema:"The auto generated employee ID" json:"id"`
	FirstName string `schema:"The employee first name" json:"first_name"`
	LastName  string `schema:"The employee last name" json:"last_name"`
	Manger Manager `schema:"The manager" json:"manger"`
	Asserts []Assert `schema:"The asserts list of the employee" json:"asserts"`
}

type ErrorResponse struct {
	at.Schema     `json:"-"`
	model.BaseResponseInfo
}

type CreateEmployeeRequest struct {
	at.RequestBody
	Employee
}

type EmployeeResponse struct {
	at.ResponseBody `json:"-"`
	at.Schema     `json:"-"`

	model.BaseResponseInfo
	Data *Employee `json:"data,omitempty" schema:"The employee data"`
}

type ListEmployeeResponse struct {
	at.ResponseBody `json:"-"`
	at.Schema     `json:"-"`
	model.BaseResponse
	Data []*Employee `json:"data,omitempty" schema:"The employee data list"`
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
func (c *employeeController) BeforeMethod(at struct{ at.BeforeMethod }, ctx context.Context) {
	log.Debug("before method")
	ctx.Next()
	return
}

// GetEmployee
func (c *employeeController) CreateEmployee(at struct {
	at.PostMapping `value:"/"`
	at.Operation   `operationId:"Create Employee" description:"This is the employee creation api"`
	at.Consumes    `values:"application/json"`
	at.Produces    `values:"application/json"`
	Parameters     struct {
		at.Parameter `name:"employee" in:"body" description:"Employee request body" `
		CreateEmployeeRequest
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a employee with ID"`
			EmployeeResponse
		}
	}
}, request *CreateEmployeeRequest) (response model.Response, err error) {
	response = new(model.BaseResponse)

	// Just for the demo purpose
	request.Employee.Id = 654321
	response.SetData(request.Employee)

	return
}


// GetEmployee
func (c *employeeController) UpdateEmployee(at struct {
	at.PutMapping `value:"/"`
	at.Operation   `operationId:"Update Employee" description:"This is the employee update api"`
	at.Consumes    `values:"application/json"`
	at.Produces    `values:"application/json"`
	Parameters     struct {
		at.Parameter `name:"employee" in:"body" description:"Employee request body" `
		CreateEmployeeRequest
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a employee with ID"`
			EmployeeResponse
		}
	}
}, request *CreateEmployeeRequest) (response model.Response, err error) {
	response = new(model.BaseResponse)

	// Just for the demo purpose
	request.Employee.Id = 67890
	response.SetData(request.Employee)

	return
}

// GetEmployee
func (c *employeeController) GetEmployee(at struct {
	at.GetMapping `value:"/{id:int}"`
	at.Operation  `operationId:"Get Employee" description:"This is get employees api"`
	at.Produces   `values:"application/json"`
	Parameters    struct {
		ID struct {
			at.Parameter `type:"integer" name:"id" in:"path" description:"Path variable employee ID" required:"true"`
		}
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a employee"`
			EmployeeResponse
		}
		StatusNotFound struct {
			at.Response `code:"404" description:"the employee you are looking for is not found"`
			ErrorResponse
		}
	}
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
		errMsg := fmt.Sprintf("employee %v is not found", id)
		response.Code = http.StatusNotFound
		response.Message = errMsg
	}

	return
}

// ListEmployee
func (c *employeeController) ListEmployee(at struct {
	at.GetMapping `value:"/"`
	at.Operation  `operationId:"List Employee" description:"This is employees list api"`
	at.Produces   `values:"application/json"`
	Responses     struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a set of employees"`
			ListEmployeeResponse
		}
		StatusNotFound struct {
			at.Response `code:"404" description:"the employees you are looking for is not found"`
			ErrorResponse
		}
	}
}) (response *ListEmployeeResponse, err error) {
	response = new(ListEmployeeResponse)
	response.Code = http.StatusOK
	response.Message = "success"
	response.Data = c.dummyData
	return
}

// DeleteEmployee
// at.DeleteEmployee is an annotation to define request mapping for http method DELETE,
func (c *employeeController) DeleteEmployee(at struct {
	at.DeleteMapping `value:"/{id:int}"`
	at.Operation     `operationId:"Delete Employee" description:"This is delete employees api"`
	at.Produces      `values:"application/json"`
	Parameters       struct {
		ID struct {
			at.Parameter `type:"integer" name:"id" in:"path" description:"Path variable employee ID" required:"true"`
		}
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns success message"`
			ErrorResponse
		}
		StatusNotFound struct {
			at.Response `code:"404" description:"the employee is not found"`
			ErrorResponse
		}
	}
}, id int) (response *EmployeeResponse, err error) {
	response = new(EmployeeResponse)
	for _, e := range c.dummyData {
		if id == e.Id {
			response.Code = http.StatusOK
			response.Message = "success"
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

// After
func (c *employeeController) AfterMethod(at struct{ at.AfterMethod }, ctx context.Context) {
	log.Debug("before method")
	ctx.Next()
	return
}

func init() {
	app.Register(
		newEmployeeController,
		swagger.OpenAPIDefinitionBuilder().
			Version("1.1.0").
			Title("HiBoot Swagger Demo Application - Simple CRUD Demo Application - 演示代码").
			Description("Simple Server is an application that demonstrate the usage of Swagger Annotations").
			Schemes("http").
			Host("localhost:8080").
			BasePath("/"),
	)
}

// Hiboot main function
func main() {
	// create new web application and run it
	web.NewApplication().
		SetProperty(app.ProfilesInclude,
			actuator.Profile,
			swagger.Profile,
			logging.Profile,
		).Run()
}
