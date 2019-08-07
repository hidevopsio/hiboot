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

// Package at provides annotations for web RestController
package at

// RestController is the annotation that declare current controller is the RESTful Controller
type RestController interface{}

// JwtRestController is the annotation that declare current controller is the RESTful Controller with JWT support
type JwtRestController interface{}

// ContextPath is the annotation that set the context path of a controller
type ContextPath interface{}

// RequestMapping is the annotation that set the RequestMapping of a controller
type RequestMapping string

// StringValue returns the value of the type RequestMapping
func (a RequestMapping) String() (value string) {
	value = string(a)
	return
}

//
//// Value set and return the value of the type RequestMapping
func (a RequestMapping) Value(str string) (value StringAnnotation) {
	value = RequestMapping(str)
	return
}

// Method is the annotation that set the RequestMethod of a controller
type Method string

// String returns the value of the type RequestMethod
func (a Method) String() (value string) {
	value = string(a)
	return
}

// Value set the value of the type RequestMethod
func (a Method) Value(str string) (value StringAnnotation) {
	value = Method(str)
	return
}

// Path is the annotation that set the Path of a controller
type Path string

// String returns the value of the type Path in string
func (a Path) String() (value string) {
	value = string(a)
	return
}

// Value set the value of the type RequestPath
func (a Path) Value(str string) (value StringAnnotation) {
	value = Path(str)
	return
}

// GetMapping is the annotation that set the GetMapping of a controller
type GetMapping struct{
	Method `value:"GET"`
}

// PostMapping is the annotation that set the PostMapping of a controller
type PostMapping struct{
	Method `value:"POST"`
}

// PutMapping is the annotation that set the PutMapping of a controller
type PutMapping struct{
	Method `value:"PUT"`
}


// DeleteMapping is the annotation that set the DeleteMapping of a controller
type DeleteMapping struct{
	Method `value:"DELETE"`
}

// AnyMapping is the annotation that set the AnyMapping of a controller
type AnyMapping struct{
	Method string `value:"ANY"`
}

// OptionsMapping is the annotation that set the OptionsMapping of a controller
type OptionsMapping struct{
	Method string `value:"OPTIONS"`
}

