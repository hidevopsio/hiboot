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

package depends

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"reflect"
	"strings"
	"fmt"
)

// ByDependency sort by the configuration dependency which specified by tag depends
type ByDependency []interface{}

func (s ByDependency) Len() int {
	return len(s)
}

func (s ByDependency) Resolve() (resolved Graph, err error) {

	var workingGraph Graph
	var node *Node
	for i, item := range s {
		// find the index of its dependency
		name := s.getFullName(item)
		dep, ok := s.findDependencies(item)
		if ok {
			node = NewNode(i, name, dep...)
		} else {
			node = NewNode(i, name)
		}
		workingGraph = append(workingGraph, node)
	}
	resolved, err = resolveGraph(workingGraph)
	return
}

func (s ByDependency) findDependencyIndex(depName string) int {
	for i, item := range s {
		name, _ := reflector.GetName(item)
		if name == depName {
			return i
		}

		pkgName := reflector.ParseObjectPkgName(item)
		fullName := pkgName + "." + name
		if fullName == depName || pkgName == depName {
			return i
		}
	}
	return -1
}

func (s ByDependency) findDependencies(item interface{}) (dep []string, ok bool) {
	// first check if contains tag depends in the embedded field
	var depName string
	depName, ok = reflector.FindEmbeddedFieldTag(item, "depends")
	if !ok {
		// else check if s is the constructor and it contains other dependencies in the input arguments
		fn := reflect.ValueOf(item)
		if fn.Type().Kind() == reflect.Func {
			numIn := fn.Type().NumIn()
			for i := 0; i < numIn; i++ {
				fnInType := fn.Type().In(i)
				depName = depName + "," + fnInType.Name()
				ok = true
			}
		}
	}

	if ok {
		dep = strings.Split(depName, ",")
		for i, dpn := range dep {
			depIdx := s.findDependencyIndex(dpn)
			if depIdx > 0 {
				depInst := s[depIdx]
				depFullName := s.getFullName(depInst)
				dep[i] = depFullName
			}
		}
	}
	return
}

func (s ByDependency) getFullName(item interface{}) (name string) {
	name, err := reflector.GetName(item)
	if err == nil {
		depPkgName := reflector.ParseObjectPkgName(item)
		name = depPkgName + "." + name
	}
	return
}

// Resolve
func Resolve(data []interface{}) (result []interface{}, err error) {

	dep := ByDependency(data)
	var resolved Graph
	resolved, err = dep.Resolve()

	if err != nil {
		log.Errorf("Failed to resolve dependencies: %s\n", err)
		displayGraph(resolved)
		fmt.Println("")
	} else {
		log.Infof("The dependency graph resolved successfully")
		for _, item := range resolved {
			result = append(result, data[item.index])
		}
	}
	return
}
