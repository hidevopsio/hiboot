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

package web

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

type FooController struct {
	Controller
}

// PUT /foo/{id}/{name}/{age}
func (c *FooController) PutByIdNameAge(id int, name string, age int) error {
	log.Debugf("FooController.Put %v %v %v", id, name, age)
	return nil
}

func TestParse(t *testing.T) {

	hdl := new(handler)

	controller := new(FooController)
	ctrlVal := reflect.ValueOf(controller)

	t.Run("should parse method with path params", func(t *testing.T) {
		method, ok := ctrlVal.Type().MethodByName("PutByIdNameAge")
		assert.Equal(t, true, ok)
		hdl.parse(method, controller, "/foo/{id}/{name}/{age}")
		log.Debug(hdl)
		assert.Equal(t, 3, len(hdl.pathParams))
		assert.Equal(t, "FooController", hdl.requests[0].typeName)
		assert.Equal(t, "int", hdl.requests[1].typeName)
		assert.Equal(t, "string", hdl.requests[2].typeName)
		assert.Equal(t, "int", hdl.requests[3].typeName)
	})
}
