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

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInvalidControllerError(t *testing.T) {
	err := InvalidControllerError{Name: "TestController"}

	assert.Equal(t, "TestController must be derived from web.Controller", err.Error())
}

func TestNotFoundError(t *testing.T) {
	err := NotFoundError{Name: "TestObject"}

	assert.Equal(t, "TestObject is not found", err.Error())
}
