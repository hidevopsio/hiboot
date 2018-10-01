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
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"reflect"
	"strings"
)

// ByDependency sort by the configuration dependency which specified by tag depends
type depResolver []*factory.MetaData

func (s depResolver) Resolve() (resolved Graph, err error) {

	var workingGraph Graph
	var node *Node
	for i, item := range s {
		// find the index of its dependency
		//name := s.getFullName(item)
		dep, ok := s.findDependencies(item)
		// TODO: temp work around
		for _, extDep := range item.ExtDep {
			workingGraph = append(workingGraph, NewNode(-1, extDep))
		}
		if ok {
			node = NewNode(i, item, dep...)
		} else {
			node = NewNode(i, item)
		}
		workingGraph = append(workingGraph, node)
	}
	resolved, err = resolveGraph(workingGraph)
	if err != nil {
		displayDependencyGraph(workingGraph, log.Error)
		displayDependencyGraph(resolved, log.Error)
	}
	return
}

func (s depResolver) findDependencyIndex(depName string) int {
	for i, item := range s {
		fullName := s.getFullName(item)

		if item.TypeName == depName {
			return i
		}

		if fullName == depName || item.PkgName == depName {
			return i
		}
	}
	return -1
}

func (s depResolver) findDependencies(item *factory.MetaData) (dep []*Node, ok bool) {
	// first check if contains tag depends in the embedded field
	var depName string
	object := item.Object
	depName, ok = reflector.FindEmbeddedFieldTag(object, "depends")
	if !ok {
		// else check if s is the constructor and it contains other dependencies in the input arguments
		outTyp, isFunc := reflector.GetFuncOutType(object)
		if isFunc {
			fn := reflect.ValueOf(object)
			numIn := fn.Type().NumIn()
			for i := 0; i < numIn; i++ {
				inTyp := fn.Type().In(i)
				indInTyp := reflector.IndirectType(inTyp)
				var name string
				for _, field := range reflector.DeepFields(outTyp) {
					indFieldTyp := reflector.IndirectType(field.Type)
					//log.Debugf("%v <> %v", indFieldTyp, indInTyp)
					if indFieldTyp == indInTyp {
						name = field.Name
						break
					}
				}
				if name == "" {
					name = reflector.GetFullNameByType(fn.Type().In(i))
				}
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
		depNames := strings.Split(depName, ",")
		for _, dp := range depNames {
			depIdx := s.findDependencyIndex(dp)
			if depIdx >= 0 {
				depInst := s[depIdx]
				dep = append(dep, NewNode(depIdx, depInst))
			} else {
				// found external dependency
				extData := &factory.MetaData{Name: dp}
				item.ExtDep = append(item.ExtDep, extData)
				dep = append(dep, NewNode(depIdx, extData))
				log.Warnf("dependency %v is not found", dp)
			}
		}
	}

	return
}

func (s depResolver) getFullName(md *factory.MetaData) (name string) {
	// check if it's func
	if md.PkgName == "" || md.TypeName == "" {
		md.PkgName, md.TypeName = reflector.GetPkgAndName(md.Object)
	}
	name = md.PkgName + "." + md.TypeName
	return
}

// Resolve resolve dependencies
func Resolve(data []*factory.MetaData) (result []*factory.MetaData, err error) {

	dep := depResolver(data)
	var resolved Graph
	resolved, err = dep.Resolve()

	if err != nil {
		log.Errorf("Failed to resolve dependencies: %s", err)
	} else {
		//log.Infof("The dependency graph resolved successfully")
		//displayDependencyGraph(resolved, log.Debug)
		for _, item := range resolved {
			if item.index >= 0 {
				result = append(result, data[item.index])
			}
		}
	}
	return
}
