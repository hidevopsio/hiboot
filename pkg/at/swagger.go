package at

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
	Swagger
}

// Api is the annotation for REST Endpoints
// e.g.
// type employeeController struct {
//   at.RestController
//   at.Api `value:"Employee Management System" description:"Operations pertaining to employee in Employee Management System"`
// }
type Api struct {
	Swagger
	Description string
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
	Swagger
	Required bool
}

// ApiResponse annotation to document other responses, in addition to the regular HTTP 200 OK, like this.
// func (c *) CreateEmployee(at struct{
//     at.PostMapping  `value:"/"`
//     at.ApiOperation `value:"Add an employee"`
//     at.ApiResponse200  `value:"Successfully retrieved list"`
//     at.ApiResponse401  `value:"You are not authorized to view the resource"`
//     at.ApiResponse403  `value:"Accessing the resource you were trying to reach is forbidden"`
//     at.ApiResponse404  `value:"The resource you were trying to reach is not found"`
//   }, request EmployeeRequest) (response Response) {
//
//   ...
// }
type ApiResponse struct {
	Swagger
}

// ApiResponse100 HTTP StatusContinue
type ApiResponse100 struct {
	ApiResponse
	Code int `value:"100"`
}

// ApiResponse101 HTTP StatusSwitchingProtocols
type ApiResponse101 struct {
	ApiResponse
	Code int `value:"101"`
}
// ApiResponse102 HTTP StatusProcessing
type ApiResponse102 struct {
	ApiResponse
	Code int `value:"102"`
}

// ApiResponse200 HTTP StatusOK
type ApiResponse200 struct {
	ApiResponse
	Code int `value:"200"`
}

// ApiResponse201 HTTP StatusCreated
type ApiResponse201 struct {
	ApiResponse
	Code int `value:"201"`
}

// ApiResponse202 HTTP StatusAccepted
type ApiResponse202 struct {
	ApiResponse
	Code int `value:"202"`
}

// ApiResponse203 HTTP StatusNonAuthoritativeInfo
type ApiResponse203 struct {
	ApiResponse
	Code int `value:"203"`
}

// ApiResponse204 HTTP StatusNoContent
type ApiResponse204 struct {
	ApiResponse
	Code int `value:"204"`
}

// ApiResponse205 HTTP StatusResetContent
type ApiResponse205 struct {
	ApiResponse
	Code int `value:"205"`
}

// ApiResponse206 HTTP StatusResetContent
type ApiResponse206 struct {
	ApiResponse
	Code int `value:"206"`
}

// ApiResponse207 HTTP StatusMultiStatus
type ApiResponse207 struct {
	ApiResponse
	Code int `value:"207"`
}

// ApiResponse208 HTTP StatusAlreadyReported
type ApiResponse208 struct {
	ApiResponse
	Code int `value:"208"`
}

// ApiResponse226 HTTP StatusIMUsed
type ApiResponse226 struct {
	ApiResponse
	Code int `value:"226"`
}

// ApiResponse300 HTTP StatusMultipleChoices
type ApiResponse300 struct {
	ApiResponse
	Code int `value:"300"`
}

// ApiResponse301 HTTP StatusMovedPermanently
type ApiResponse301 struct {
	ApiResponse
	Code int `value:"301"`
}

// ApiResponse302 HTTP StatusFound
type ApiResponse302 struct {
	ApiResponse
	Code int `value:"302"`
}

// ApiResponse303 HTTP StatusSeeOther
type ApiResponse303 struct {
	ApiResponse
	Code int `value:"303"`
}

// ApiResponse304 HTTP StatusNotModified
type ApiResponse304 struct {
	ApiResponse
	Code int `value:"304"`
}

// ApiResponse305 HTTP StatusUseProxy
type ApiResponse305 struct {
	ApiResponse
	Code int `value:"305"`
}

// ApiResponse307 HTTP StatusTemporaryRedirect
type ApiResponse307 struct {
	ApiResponse
	Code int `value:"307"`
}

// ApiResponse308 HTTP StatusPermanentRedirect
type ApiResponse308 struct {
	ApiResponse
	Code int `value:"308"`
}

// ApiResponse400 HTTP StatusBadRequest
type ApiResponse400 struct {
	ApiResponse
	Code int `value:"400"`
}

// ApiResponse401 HTTP StatusUnauthorized
type ApiResponse401 struct {
	ApiResponse
	Code int `value:"401"`
}

// ApiResponse402 HTTP StatusPaymentRequired
type ApiResponse402 struct {
	ApiResponse
	Code int `value:"402"`
}

// ApiResponse403 HTTP StatusForbidden
type ApiResponse403 struct {
	ApiResponse
	Code int `value:"403"`
}

// ApiResponse404 HTTP StatusNotFound
type ApiResponse404 struct {
	ApiResponse
	Code int `value:"404"`
}

// ApiResponse405 HTTP StatusMethodNotAllowed
type ApiResponse405 struct {
	ApiResponse
	Code int `value:"405"`
}

// ApiResponse406 HTTP StatusNotAcceptable
type ApiResponse406 struct {
	ApiResponse
	Code int `value:"406"`
}

// ApiResponse407 HTTP StatusProxyAuthRequired
type ApiResponse407 struct {
	ApiResponse
	Code int `value:"407"`
}

// ApiResponse408 HTTP StatusRequestTimeout
type ApiResponse408 struct {
	ApiResponse
	Code int `value:"408"`
}

// ApiResponse409 HTTP StatusConflict
type ApiResponse409 struct {
	ApiResponse
	Code int `value:"409"`
}

// ApiResponse410 HTTP StatusGone
type ApiResponse410 struct {
	ApiResponse
	Code int `value:"410"`
}

// ApiResponse411 HTTP StatusLengthRequired
type ApiResponse411 struct {
	ApiResponse
	Code int `value:"411"`
}

// ApiResponse412 HTTP StatusPreconditionFailed
type ApiResponse412 struct {
	ApiResponse
	Code int `value:"412"`
}

// ApiResponse413 HTTP StatusRequestEntityTooLarge
type ApiResponse413 struct {
	ApiResponse
	Code int `value:"413"`
}

// ApiResponse414 HTTP StatusRequestURITooLong
type ApiResponse414 struct {
	ApiResponse
	Code int `value:"414"`
}

// ApiResponse415 HTTP StatusUnsupportedMediaType
type ApiResponse415 struct {
	ApiResponse
	Code int `value:"415"`
}

// ApiResponse416 HTTP StatusRequestedRangeNotSatisfiable
type ApiResponse416 struct {
	ApiResponse
	Code int `value:"416"`
}

// ApiResponse417 HTTP StatusExpectationFailed
type ApiResponse417 struct {
	ApiResponse
	Code int `value:"417"`
}

// ApiResponse418 HTTP StatusTeapot
type ApiResponse418 struct {
	ApiResponse
	Code int `value:"418"`
}

// ApiResponse421 HTTP StatusMisdirectedRequest
type ApiResponse421 struct {
	ApiResponse
	Code int `value:"421"`
}

// ApiResponse422 HTTP StatusUnprocessableEntity
type ApiResponse422 struct {
	ApiResponse
	Code int `value:"422"`
}

// ApiResponse423 HTTP StatusLocked
type ApiResponse423 struct {
	ApiResponse
	Code int `value:"423"`
}

// ApiResponse424 HTTP StatusFailedDependency
type ApiResponse424 struct {
	ApiResponse
	Code int `value:"424"`
}

// ApiResponse425 HTTP StatusTooEarly
type ApiResponse425 struct {
	ApiResponse
	Code int `value:"425"`
}

// ApiResponse426 HTTP StatusUpgradeRequired
type ApiResponse426 struct {
	ApiResponse
	Code int `value:"426"`
}

// ApiResponse428 HTTP StatusPreconditionRequired
type ApiResponse428 struct {
	ApiResponse
	Code int `value:"428"`
}

// ApiResponse429 HTTP StatusTooManyRequests
type ApiResponse429 struct {
	ApiResponse
	Code int `value:"429"`
}

// ApiResponse431 HTTP StatusRequestHeaderFieldsTooLarge
type ApiResponse431 struct {
	ApiResponse
	Code int `value:"431"`
}

// ApiResponse451 HTTP StatusUnavailableForLegalReasons
type ApiResponse451 struct {
	ApiResponse
	Code int `value:"451"`
}

// ApiResponse500 HTTP StatusInternalServerError
type ApiResponse500 struct {
	ApiResponse
	Code int `value:"500"`
}

// ApiResponse501 HTTP StatusNotImplemented
type ApiResponse501 struct {
	ApiResponse
	Code int `value:"501"`
}

// ApiResponse502 HTTP StatusBadGateway
type ApiResponse502 struct {
	ApiResponse
	Code int `value:"502"`
}

// ApiResponse503 HTTP StatusServiceUnavailable
type ApiResponse503 struct {
	ApiResponse
	Code int `value:"503"`
}

// ApiResponse504 HTTP StatusGatewayTimeout
type ApiResponse504 struct {
	ApiResponse
	Code int `value:"504"`
}

// ApiResponse505 HTTP StatusHTTPVersionNotSupported
type ApiResponse505 struct {
	ApiResponse
	Code int `value:"505"`
}

// ApiResponse506 HTTP StatusVariantAlsoNegotiates
type ApiResponse506 struct {
	ApiResponse
	Code int `value:"506"`
}

// ApiResponse507 HTTP StatusInsufficientStorage
type ApiResponse507 struct {
	ApiResponse
	Code int `value:"507"`
}

// ApiResponse508 HTTP StatusLoopDetected
type ApiResponse508 struct {
	ApiResponse
	Code int `value:"508"`
}

// ApiResponse510 HTTP StatusNotExtended
type ApiResponse510 struct {
	ApiResponse
	Code int `value:"510"`
}

// ApiResponse511 HTTP StatusNetworkAuthenticationRequired
type ApiResponse511 struct {
	ApiResponse
	Code int `value:"511"`
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