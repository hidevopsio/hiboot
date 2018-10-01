// Copyright (c) 2016 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer
//    in this position and unchanged.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// this file is copied from https://github.com/dnaeon/go-dependency-graph-algorithm

package depends

import (
	"errors"
	"fmt"

	"github.com/deckarep/golang-set"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/log"
)

// Node represents a single node in the graph with it's dependencies
type Node struct {
	// Index of the node
	index int

	// Name of the node
	data *factory.MetaData

	// Dependencies of the node
	deps []*Node
}

// ErrCircularDependency report that circular dependency found
var ErrCircularDependency = errors.New("circular dependency found")

// NewNode creates a new node
func NewNode(index int, data interface{}, deps ...*Node) *Node {
	var md *factory.MetaData
	if reflect.TypeOf(data).Kind() == reflect.String {
		md = &factory.MetaData{Name: data.(string)}
	} else {
		md = data.(*factory.MetaData)
	}

	n := &Node{
		index: index,
		data:  md,
		deps:  deps,
	}

	return n
}

type Graph []*Node

// Resolves the dependency graph
func resolveGraph(graph Graph) (Graph, error) {
	// A map containing the node names and the actual node object
	nodeNames := make(map[string]*Node)

	// A map containing the nodes and their dependencies
	nodeDependencies := make(map[string]mapset.Set)

	// Populate the maps
	for _, node := range graph {
		nodeNames[node.data.Name] = node

		dependencySet := mapset.NewSet()
		for _, dep := range node.deps {
			dependencySet.Add(dep)
			log.Info(dep)
		}
		nodeDependencies[node.data.Name] = dependencySet
	}

	log.Info(nodeDependencies)

	// Iteratively find and remove nodes from the graph which have no dependencies.
	// If at some point there are still nodes in the graph and we cannot find
	// nodes without dependencies, that means we have a circular dependency
	var resolved Graph
	loop := 0
	for len(nodeDependencies) != 0 {
		loop++
		// Get all nodes from the graph which have no dependencies
		readySet := mapset.NewSet()
		for name, deps := range nodeDependencies {
			if deps.Cardinality() == 0 {
				readySet.Add(name)
			}
		}

		// If there aren't any ready nodes, then we have a circular dependency
		if readySet.Cardinality() == 0 {
			var g Graph
			for name := range nodeDependencies {
				g = append(g, nodeNames[name])
			}

			return g, ErrCircularDependency
		}

		// Remove the ready nodes and add them to the resolved graph
		for name := range readySet.Iter() {
			delete(nodeDependencies, name.(string))
			resolved = append(resolved, nodeNames[name.(string)])
		}

		// Also make sure to remove the ready nodes from the
		// remaining node dependencies as well
		for name, deps := range nodeDependencies {
			diff := deps.Difference(readySet)
			nodeDependencies[name] = diff
		}
	}
	log.Debugf("loop: %d", loop)

	return resolved, nil
}

// Displays the dependency graph
func displayDependencyGraph(graph Graph, logger func(v ...interface{})) {
	output := "\n\nDependency tree:\n"
	for i, node := range graph {
		if len(node.deps) == 0 {
			output += fmt.Sprintf("    %d(%d): %s ->\n", i, node.index, node.data.Name)
		} else {
			for _, dep := range node.deps {
				output += fmt.Sprintf("    %d(%d): %s -> %s\n", i, node.index, node.data.Name, dep.data.Name)
			}
		}
	}
	if reflect.TypeOf(logger).Kind() == reflect.Func {
		logger(output)
	}
}
