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
	"github.com/hidevopsio/hiboot/pkg/system/types"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/hidevopsio/hiboot/pkg/utils/str"
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
		dep, ok := s.findDependencies(item)
		//// TODO: temp work around
		//for _, extDep := range item.ExtDep {
		//	workingGraph = append(workingGraph, NewNode(-1, extDep))
		//}
		if ok {
			node = NewNode(i, item, dep...)
		} else {
			node = NewNode(i, item)
		}
		workingGraph = append(workingGraph, node)
	}
	resolved, err = resolveGraph(workingGraph)
	displayDependencyGraph("working graph", workingGraph, log.Debug)
	return
}

func (s depResolver) findDependencyIndex(depName string) int {
	for i, item := range s {
		if item.TypeName == depName {
			return i
		}

		// find item name or package name
		if item.ShortName == depName || item.Name == depName || item.PkgName == depName {
			return i
		}

		// else find method name
		if item.Kind == types.Method {
			method := item.Object.(reflect.Method)
			methodName := str.ToLowerCamel(method.Name)
			if depName == methodName || depName == item.PkgName+"."+methodName {
				return i
			}
		}
	}
	return -1
}

func (s depResolver) findDependencies(item *factory.MetaData) (dep []*Node, ok bool) {
	// first check if contains tag depends in the embedded field
	var depName string
	object := item.Object
	switch item.Kind {
	case types.Func:
		fn := reflect.ValueOf(object)
		numIn := fn.Type().NumIn()
		for i := 0; i < numIn; i++ {
			inTyp := fn.Type().In(i)
			indInTyp := reflector.IndirectType(inTyp)
			var name string
			for _, field := range reflector.DeepFields(item.Type) {
				indFieldTyp := reflector.IndirectType(field.Type)
				//log.Debugf("%v <> %v", indFieldTyp, indInTyp)
				if indFieldTyp == indInTyp {
					name = str.ToLowerCamel(field.Name)
					depPkgName := io.DirName(field.Type.PkgPath())
					if depPkgName != "" {
						name = depPkgName + "." + name
					}
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
	case types.Method:
		// TODO: too many duplicated code, optimize it
		method := object.(reflect.Method)
		numIn := method.Type.NumIn()
		for i := 1; i < numIn; i++ {
			inTyp := method.Type.In(i)
			indInTyp := reflector.IndirectType(inTyp)
			var name string
			if item.Type != nil {
				for _, field := range reflector.DeepFields(item.Type) {
					indFieldTyp := reflector.IndirectType(field.Type)
					//log.Debugf("%v <> %v", indFieldTyp, indInTyp)
					if indFieldTyp == indInTyp {
						name = str.ToLowerCamel(field.Name)
						depPkgName := io.DirName(indFieldTyp.PkgPath())
						if depPkgName != "" {
							name = depPkgName + "." + name
						}
						break
					}
				}
			}
			if name == "" {
				name = reflector.GetFullNameByType(method.Type.In(i))
			}
			if depName == "" {
				depName = name
			} else {
				depName = depName + "," + name
			}
			ok = true
		}
	default:
		depName, ok = reflector.FindEmbeddedFieldTag(object, "depends")
	}

	// iterate dependencies
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

// Resolve resolve dependencies
func Resolve(data []*factory.MetaData) (result []*factory.MetaData, err error) {
	if len(data) != 0 {

		dep := depResolver(data)
		var resolved Graph
		resolved, err = dep.Resolve()

		if err != nil {
			log.Errorf("Failed to resolve dependencies: %s", err)
			displayDependencyGraph("broken graph", resolved, log.Error)
		} else {
			log.Infof("The dependency graph resolved successfully")
			displayDependencyGraph("resolved graph", resolved, log.Debug)
			for _, item := range resolved {
				if item.index >= 0 {
					result = append(result, data[item.index])
				}
			}
		}
	}

	return
}
