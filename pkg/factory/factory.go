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

// Package factory provides InstantiateFactory and ConfigurableFactory interface
package factory

import (
	"github.com/hidevopsio/hiboot/pkg/system"
	"reflect"
)

type Factory interface{}

type InstantiateFactory interface {
	Initialized() bool
	SetInstance(name string, instance interface{}) (err error)
	GetInstance(name string) (inst interface{})
	Items() map[string]interface{}
	AppendComponent(c ...interface{})
}

type ConfigurableFactory interface {
	InstantiateFactory
	SystemConfiguration() *system.Configuration
	Configuration(name string) interface{}
}

type MetaData struct {
	Kind     reflect.Kind
	Name     string
	TypeName string
	PkgName  string
	Object   interface{}
	ExtDep   []string
}
