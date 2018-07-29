package model

var (
	RequestTypeBody = "body"
	RequestTypeParams = "params"
	RequestTypeForm = "form"
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
