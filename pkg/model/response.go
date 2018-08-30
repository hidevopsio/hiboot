// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

type Response interface {
	SetCode(code int)
	GetCode() int
	SetMessage(message string)
	GetMessage() string
	SetData(data interface{})
	GetData() interface{}
}

type BaseResponse struct {
	Response
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r *BaseResponse) SetCode(code int) {
	r.Code = code
}

func (r *BaseResponse) GetCode() int {
	return r.Code
}

func (r *BaseResponse) SetMessage(message string) {
	r.Message = message
}

func (r *BaseResponse) GetMessage() string {
	return r.Message
}

func (r *BaseResponse) SetData(data interface{}) {
	r.Data = data
}

func (r *BaseResponse) GetData() interface{} {
	return r.Data
}
