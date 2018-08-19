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

package str

import (
	"unicode"
	"strings"
)

const EmptyString  = ""

// UpperFirst upper case first character of specific string
func UpperFirst(str string) string {
	return strings.Title(str)
}

// LowerFirst lower case first character of specific string
func LowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return EmptyString
}


// StringInSlice check if specific string is in slice
func InSlice(a string, list []string) bool {

	var retVal bool

	for _, b := range list {
		if b == a {
			retVal = true
			break
		}
	}
	return retVal
}