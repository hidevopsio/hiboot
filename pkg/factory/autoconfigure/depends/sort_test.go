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

package depends_test

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/factory/autoconfigure/depends"
	"github.com/hidevopsio/hiboot/pkg/factory/autoconfigure/depends/fake"
	"github.com/hidevopsio/hiboot/pkg/factory/autoconfigure/depends/foo"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"sort"
	"testing"
)

type fooConfiguration struct {
	app.Configuration
}

type barConfiguration struct {
	app.Configuration
}

type childConfiguration struct {
	app.Configuration `depends:"parentConfiguration"`
}

type parentConfiguration struct {
	app.Configuration `depends:"grantConfiguration"`
}

type grantConfiguration struct {
	app.Configuration `depends:"fake.Configuration"`
}

type cyclingChildConfiguration struct {
	app.Configuration `depends:"cyclingParentConfiguration"`
}

type cyclingParentConfiguration struct {
	app.Configuration `depends:"cyclingGrantConfiguration"`
}

type cyclingGrantConfiguration struct {
	app.Configuration `depends:"cyclingChildConfiguration"`
}

func TestSort(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	testData := []struct {
		title          string
		configurations []interface{}
		err            error
	}{
		{
			title:          "should sort dependencies",
			configurations: []interface{}{new(fooConfiguration), new(childConfiguration), new(parentConfiguration), new(grantConfiguration), new(fake.Configuration), new(foo.Configuration), new(barConfiguration)},
			err:            nil,
		},
		{
			title:          "should fail to sort with cycling dependencies",
			configurations: []interface{}{new(cyclingChildConfiguration), new(cyclingParentConfiguration), new(cyclingGrantConfiguration)},
			err:            nil,
		},
	}

	for _, data := range testData {
		t.Run(data.title, func(t *testing.T) {
			sort.Sort(depends.ByDependency(data.configurations))
			for i, item := range data.configurations {
				name, _ := reflector.GetName(item)
				pkgName := reflector.ParseObjectPkgName(item)
				t, _ := reflector.FindEmbeddedFieldTag(item, "depends")
				log.Debugf("[after] %v: %v.%v -> %v", i, pkgName, name, t)
			}
		})
	}
}
