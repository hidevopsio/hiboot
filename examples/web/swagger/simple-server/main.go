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
	"time"
)

type Asset struct {
	at.Schema
	ID         int    `schema:"The asset ID" json:"id"`
	Name       string `schema:"The asset name" json:"name"`
	Amount     float64
	Type       string `schema:"The asset type"`
	ExpirationTime time.Time
}

type AddAssetsResponse struct {
	at.ResponseBody `json:"-"`
	at.Schema

	model.BaseResponseInfo
	Data []*Asset `json:"data,omitempty" schema:"The employee data"`
}

type Manager struct {
	at.Schema
	ID int `schema:"The manager ID" json:"id"`
	Name string `schema:"The manager name of the employee" json:"name"`
}

type Employee struct {
	at.Schema
	Id        int     `schema:"The auto generated employee ID" json:"id"`
	FirstName string  `schema:"The employee first name" json:"first_name"`
	LastName  string  `schema:"The employee last name" json:"last_name"`
	Manger    Manager `schema:"The manager" json:"manger"`
	Assets    []Asset `schema:"The assets list of the employee" json:"assets"`
}

type ErrorResponse struct {
	at.Schema
	model.BaseResponseInfo
}


type UpdateEmployeeRequest struct {
	at.RequestBody
	Employee
}

type CreateEmployeeRequest struct {
	at.RequestBody
	Employee
}

type EmployeeResponse struct {
	at.ResponseBody `json:"-"`
	at.Schema

	model.BaseResponseInfo
	Data *Employee `json:"data,omitempty" schema:"The employee data"`
}

type ListEmployeeResponse struct {
	at.ResponseBody `json:"-"`
	at.Schema
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
	log.Debugf("%v %v before method", ctx.GetCurrentRoute().Method(), ctx.GetCurrentRoute().Path())
	ctx.Next()
	return
}

// GetEmployee
func (c *employeeController) CreateEmployee(at struct {
	at.PostMapping `value:"/"`
	at.Operation   `id:"Create Employee" description:"This is the employee creation api"`
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
}, request *CreateEmployeeRequest) (response *EmployeeResponse, err error) {
	response = new(EmployeeResponse)

	// Just for the demo purpose
	request.Employee.Id = 654321
	response.SetCode(http.StatusOK)
	response.Data = &request.Employee

	return
}


// GetEmployee
func (c *employeeController) UpdateEmployee(at struct {
	at.PutMapping `value:"/"`
	at.Operation   `id:"Update Employee" description:"This is the employee update api"`
	at.Consumes    `values:"application/json"`
	at.Produces    `values:"application/json"`
	Parameters     struct {
		at.Parameter `name:"employee" in:"body" description:"Employee request body" `
		UpdateEmployeeRequest
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a employee with ID"`
			EmployeeResponse
		}
	}
}, request *UpdateEmployeeRequest) (response model.Response, err error) {
	response = new(model.BaseResponse)

	// Just for the demo purpose
	request.Employee.Id = 67890
	response.SetData(request.Employee)

	return
}

// GetEmployee
func (c *employeeController) GetEmployee(at struct {
	at.GetMapping `value:"/{id}"`
	at.Operation  `id:"Get Employee" description:"This is get employees api"`
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


// GetEmployeeName
func (c *employeeController) GetEmployeeName(at struct {
	at.GetMapping `value:"/{id}/name"`
	at.Operation  `id:"Get Employee Name" description:"This is the api that get employee name"`
	at.Produces   `values:"text/plain"`
	Parameters    struct {
		ID struct {
			at.Parameter `type:"integer" name:"id" in:"path" description:"Path variable employee ID" required:"true"`
		}
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns the employee name"`
			at.Schema   `type:"string" description:"contains the actual employee name as plain text"`
		}
		StatusNotFound struct {
			at.Response `code:"404" description:"employee is not found"`
			at.Schema   `type:"string" description:"Report 'not found' error message"`
		}
	}
}, id int) (name string) {
	return "Donald Trump"
}


// ListEmployee
func (c *employeeController) ListEmployee(at struct {
	at.GetMapping `value:"/"`
	at.Operation  `id:"List Employee" description:"This is employees list api"`
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
	at.DeleteMapping `value:"/{id}"`
	at.Operation     `id:"Delete Employee" description:"This is delete employees api"`
	at.Produces      `values:"application/json"`
	Parameters       struct {
		at.Parameter `type:"integer" name:"id" in:"path" description:"Path variable employee ID" required:"true"`
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
}, id int) (response model.ResponseInfo, err error) {
	response = new(model.BaseResponseInfo)
	//for _, e := range c.dummyData {
	//	if id == e.Id {
	//		break
	//	}
	//}
	return
}

// AddEmployeeAsserts
func (c *employeeController) AddEmployeeAsserts(at struct {
	at.PostMapping `value:"/add-assets"`
	at.Operation     `id:"Add Employee's Assets" description:"This is the api that adding assets for employees"`
	at.Consumes      `values:"application/json"`
	at.Produces      `values:"application/json"`
	Parameters    struct {
		at.Parameter `in:"body" description:"Employee request body" `
		at.Schema `value:"array"`
		Assets []*Asset
	}
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"returns a employee with ID"`
			at.Schema `value:"array"`
			Assets []*Asset
		}
	}
}, assets []*Asset) (response model.Response, err error) {
	response = new(model.BaseResponse)
	response.SetData(assets)
	return
}

// After
func (c *employeeController) AfterMethod(at struct{ at.AfterMethod }, ctx context.Context) {
	log.Debugf("%v %v after method", ctx.GetCurrentRoute().Method(), ctx.GetCurrentRoute().Path())
	ctx.Next()
	return
}

func init() {
	app.Register(
		newEmployeeController,
		swagger.ApiInfoBuilder().
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
