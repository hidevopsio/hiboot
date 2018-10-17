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

// Package sort provides utility that sort slice by length
package sort

import (
	"sort"
)

type byLen []string

// get slice length
func (a byLen) Len() int {
	return len(a)
}

// Less check which element is less
func (a byLen) Less(i, j int) bool {
	return len(a[i]) < len(a[j])
}

// Swap swap elements
func (a byLen) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// ByLen sort by length
func ByLen(s []string) {

	sort.Sort(byLen(s))

}
