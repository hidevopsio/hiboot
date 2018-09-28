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

package depends

import (
	"fmt"
	"testing"
	"github.com/hidevopsio/hiboot/pkg/log"
)

func TestDepTree(t *testing.T) {
	//
	// A working dependency graph
	//
	nodeA := NewNode(0, "A")
	nodeB := NewNode(1, "B")
	nodeC := NewNode(2, "C", "A")
	nodeD := NewNode(3, "D", "B")
	nodeE := NewNode(4, "E", "C", "D")
	nodeF := NewNode(5, "F", "A", "B")
	nodeG := NewNode(6, "G", "E", "F")
	nodeH := NewNode(7, "H", "G")
	nodeI := NewNode(8, "I", "A")
	nodeJ := NewNode(8, "J", "B")
	nodeK := NewNode(10, "K")

	var workingGraph Graph
	workingGraph = append(workingGraph, nodeA, nodeB, nodeC, nodeD, nodeE, nodeF, nodeG, nodeH, nodeI, nodeJ, nodeK)

	fmt.Printf(">>> A working dependency graph\n")
	displayGraph(workingGraph)

	resolved, err := resolveGraph(workingGraph)
	if err != nil {
		fmt.Printf("Failed to resolve dependency graph: %s\n", err)
	} else {
		fmt.Println("The dependency graph resolved successfully")
	}

	for _, node := range resolved {
		fmt.Println(node.name)
	}

	//
	// A broken dependency graph with circular dependency
	//
	nodeA = NewNode(11, "A", "I")

	var brokenGraph Graph
	brokenGraph = append(brokenGraph, nodeA, nodeB, nodeC, nodeD, nodeE, nodeF, nodeG, nodeH, nodeI, nodeJ, nodeK)

	fmt.Printf(">>> A broken dependency graph with circular dependency\n")
	displayGraph(brokenGraph)

	resolved, err = resolveGraph(brokenGraph)
	if err != nil {
		fmt.Printf("Failed to resolve dependency graph: %s\n", err)
	} else {
		fmt.Println("The dependency graph resolved successfully")
	}
}

func TestDep(t *testing.T) {
	var workingGraph Graph
	workingGraph = append(workingGraph,
		NewNode(0, "a", "b"),
		NewNode(1, "b", "c"),
		NewNode(2, "c", "d", "e"),
		NewNode(3, "d"),
		NewNode(4, "e"))

	fmt.Printf(">>> A working dependency graph\n")
	displayGraph(workingGraph)

	resolved, err := resolveGraph(workingGraph)
	if err != nil {
		fmt.Printf("Failed to resolve dependency graph: %s\n", err)
	} else {
		fmt.Println("The dependency graph resolved successfully")
	}
	log.Debugf("resolved: %v", resolved)
}