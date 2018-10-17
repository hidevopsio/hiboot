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

// Response is the interface of rest controller's Response
type Response interface {
	// Set error code
	SetCode(code int)
	// Get error code
	GetCode() int
	// Set message
	SetMessage(message string)
	// Get message
	GetMessage() string
	// Set data, the data will be serialized to json string
	SetData(data interface{})
	// Get data
	GetData() interface{}
}

// BaseResponse is the implementation of rest controller's Response
type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// SetCode set error code
func (r *BaseResponse) SetCode(code int) {
	r.Code = code
}

// GetCode get error code
func (r *BaseResponse) GetCode() int {
	return r.Code
}

// SetMessage set message
func (r *BaseResponse) SetMessage(message string) {
	r.Message = message
}

// GetMessage get message
func (r *BaseResponse) GetMessage() string {
	return r.Message
}

// SetData the data will be serialized to json string
func (r *BaseResponse) SetData(data interface{}) {
	r.Data = data
}

// GetData get data
func (r *BaseResponse) GetData() interface{} {
	return r.Data
}
