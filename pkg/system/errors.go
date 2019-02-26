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

package system

import "fmt"

// ErrInvalidController invalid controller
type ErrInvalidController struct {
	Name string
}

func (e *ErrInvalidController) Error() string {
	// TODO: locale
	return fmt.Sprintf("%v must be derived from at.RestController", e.Name)
}

// ErrNotFound resource not found error
type ErrNotFound struct {
	Name string
}

func (e *ErrNotFound) Error() string {
	// TODO: locale
	return fmt.Sprintf("%v is not found", e.Name)
}
