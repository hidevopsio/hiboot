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
type RestController struct {
	Annotation
}

// JwtRestController is the annotation that declare current controller is the RESTful Controller with JWT support
type JwtRestController struct {
	RestController
}

// Middleware is the annotation that annotate the controller or method use middleware
type Middleware struct {
	Annotation
}

// UseJwt is the annotation that that annotate the controller or method use JWT middleware
type UseJwt struct {
	Middleware
}

// ContextPath is the annotation that set the context path of a controller
type ContextPath struct {
	Annotation
}

// RequestMapping is the annotation that set the RequestMapping of a controller
type RequestMapping struct {
	HttpMethod
}

// Method is the annotation that set the RequestMethod of a controller
type Method struct {
	Annotation
}

// HttpMethod is the annotation that the http method of a controller
type HttpMethod struct {
	Annotation
}

// BeforeMethod is the annotation that set the called before the http method of a controller
type BeforeMethod struct {
	HttpMethod
}

// AfterMethod is the annotation that set the called after the http method of a controller
type AfterMethod struct {
	HttpMethod
}

// GetMapping is the annotation that set the GetMapping of a controller
type GetMapping struct {
	RequestMapping
	Method string `value:"GET" json:"-"`
}

// PostMapping is the annotation that set the PostMapping of a controller
type PostMapping struct {
	RequestMapping
	Method string  `value:"POST" json:"-"`
}

// PutMapping is the annotation that set the PutMapping of a controller
type PutMapping struct {
	RequestMapping
	Method string  `value:"PUT" json:"-"`
}

// PatchMapping is the annotation that set the PatchMapping of a controller
type PatchMapping struct {
	RequestMapping
	Method string  `value:"PATCH" json:"-"`
}

// DeleteMapping is the annotation that set the DeleteMapping of a controller
type DeleteMapping struct {
	RequestMapping
	Method string  `value:"DELETE" json:"-"`
}

// AnyMapping is the annotation that set the AnyMapping of a controller
type AnyMapping struct {
	RequestMapping
	Method string  `value:"ANY" json:"-"`
}

// OptionsMapping is the annotation that set the OptionsMapping of a controller
type OptionsMapping struct {
	RequestMapping
	Method string  `value:"OPTIONS" json:"-"`
}

// TraceMapping is the annotation that set the TraceMapping of a controller
type TraceMapping struct {
	RequestMapping
	Method string  `value:"TRACE" json:"-"`
}