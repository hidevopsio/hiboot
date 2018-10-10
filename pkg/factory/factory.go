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
)

const (
	InstantiateFactoryName  = "factory.instantiateFactory"
	ConfigurableFactoryName = "factory.configurableFactory"
)

type Factory interface{}

// InstantiateFactory instantiate factory interface
type InstantiateFactory interface {
	Initialized() bool
	SetInstance(params ...interface{}) (err error)
	GetInstance(params ...interface{}) (retVal interface{})
	GetInstances(name string) (retVal []interface{})
	Items() map[string]interface{}
	AppendComponent(c ...interface{})
}

// ConfigurableFactory configurable factory interface
type ConfigurableFactory interface {
	InstantiateFactory
	SystemConfiguration() *system.Configuration
	Configuration(name string) interface{}
}

// Configuration configuration interface
type Configuration interface {
	dependencies(name string) (deps []string)
	setDependencies(name string, value []string)
}

type depsMap map[string][]string

type Deps struct {
	deps depsMap
}

func (c *Deps) ensure() {
	if c.deps == nil {
		c.deps = make(depsMap)
	}
}

func (c *Deps) Get(name string) (deps []string) {
	c.ensure()

	deps = c.deps[name]

	return
}

func (c *Deps) Set(name string, value []string) {
	c.ensure()

	c.deps[name] = value

	return
}
