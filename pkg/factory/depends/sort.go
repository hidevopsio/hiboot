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
	"github.com/hidevopsio/hiboot/pkg/factory"
)

// ByDependency sort by the configuration dependency which specified by tag depends
type Dep []*factory.MetaData

func (s Dep) Resolve() (resolved Graph, err error) {

	var workingGraph Graph
	var node *Node
	for i, item := range s {
		// find the index of its dependency
		name := s.getFullName(item)
		dep, ok := s.findDependencies(item)
		for _, extDep := range item.ExtDep {
			workingGraph = append(workingGraph, NewNode(-1, extDep))
		}
		if ok {
			node = NewNode(i, name, dep...)
		} else {
			node = NewNode(i, name)
		}
		workingGraph = append(workingGraph, node)
	}
	resolved, err = resolveGraph(workingGraph)
	if err != nil {
		displayDependencyGraph(workingGraph, log.Error)
	}
	return
}

func (s Dep) findDependencyIndex(depName string) int {
	for i, item := range s {
		fullName := s.getFullName(item)

		if item.Name == depName {
			return i
		}

		if fullName == depName || item.PkgName == depName {
			return i
		}
	}
	return -1
}

func (s Dep) findDependencies(item *factory.MetaData) (dep []string, ok bool) {
	// first check if contains tag depends in the embedded field
	var depName string
	depName, ok = reflector.FindEmbeddedFieldTag(item.Object, "depends")
	if !ok {
		// else check if s is the constructor and it contains other dependencies in the input arguments
		fn := reflect.ValueOf(item.Object)
		if item.Kind == reflect.Func {
			numIn := fn.Type().NumIn()
			for i := 0; i < numIn; i++ {
				name := reflector.GetFullNameByType(fn.Type().In(i))
				if depName == "" {
					depName = name
				} else {
					depName = depName + "," + name
				}
				ok = true
			}
		}
	}

	if ok {
		dep = strings.Split(depName, ",")
		for i, dp := range dep {
			depIdx := s.findDependencyIndex(dp)
			if depIdx >= 0 {
				depInst := s[depIdx]
				depFullName := s.getFullName(depInst)
				dep[i] = depFullName
			} else {
				// found external dependency
				item.ExtDep = append(item.ExtDep, dp)
				dep[i] = dp
				log.Warnf("dependency %v is not found", dp)
			}
		}
	}
	return
}

func (s Dep) getFullName(md *factory.MetaData) (name string) {
	// check if it's func
	if md.PkgName == "" || md.Name == "" {
		md.PkgName, md.Name = reflector.GetPkgAndName(md.Object)
	}
	name = md.PkgName + "." + md.Name
	return
}

// TODO: support multi-dimensional array/slice
// Resolve
func Resolve(data []*factory.MetaData) (result []*factory.MetaData, err error) {

	dep := Dep(data)
	var resolved Graph
	resolved, err = dep.Resolve()

	if err != nil {
		log.Errorf("Failed to resolve dependencies: %s", err)
		displayDependencyGraph(resolved, log.Error)
	} else {
		//log.Infof("The dependency graph resolved successfully")
		displayDependencyGraph(resolved, log.Debug)
		for _, item := range resolved {
			if item.index >= 0 {
				result = append(result, data[item.index])
			}
		}
	}
	return
}
