package model

var (
	RequestTypeBody = "RequestBody"
	RequestTypeParams = "RequestParams"
	RequestTypeForm = "RequestForm"
	Context = "Context"
)

type RequestBody struct {
}

func (r *RequestBody) RequestType() string  {
	return RequestTypeBody
}

type RequestForm struct {
}

func (r *RequestForm) RequestType() string  {
	return RequestTypeForm
}

type RequestParams struct {
}

func (r *RequestParams) RequestType() string  {
	return RequestTypeParams
}
