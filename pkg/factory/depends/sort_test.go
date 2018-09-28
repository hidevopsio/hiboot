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
	"github.com/hidevopsio/hiboot/pkg/log"
	"testing"
	"github.com/hidevopsio/hiboot/pkg/factory/depends"
	"github.com/hidevopsio/hiboot/pkg/factory/depends/bar"
	"github.com/hidevopsio/hiboot/pkg/factory/depends/fake"
	"github.com/hidevopsio/hiboot/pkg/factory/depends/foo"
	"github.com/magiconair/properties/assert"
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

type circularChildConfiguration struct {
	app.Configuration `depends:"cyclingParentConfiguration"`
}

type circularParentConfiguration struct {
	app.Configuration `depends:"cyclingGrantConfiguration"`
}

type circularGrantConfiguration struct {
	app.Configuration `depends:"cyclingParentConfiguration"`
}

type circularChildConfiguration2 struct {
	app.Configuration `depends:"cyclingParentConfiguration2"`
}

type circularParentConfiguration2 struct {
	app.Configuration `depends:"cyclingGrantConfiguration2"`
}

type circularGrantConfiguration2 struct {
	app.Configuration `depends:"cyclingChildConfiguration2"`
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
			configurations: []interface{}{new(fooConfiguration), new(bar.Configuration), new(childConfiguration), new(parentConfiguration), new(grantConfiguration), new(fake.Configuration), new(foo.Configuration), new(barConfiguration)},
			err:            nil,
		},
		{
			title:          "should fail to sort with cycling dependencies",
			configurations: []interface{}{new(circularChildConfiguration), new(circularParentConfiguration), new(circularGrantConfiguration)},
			err:            depends.ErrCircularDependency,
		},
		{
			title:          "should fail to sort with cycling dependencies 2",
			configurations: []interface{}{new(circularChildConfiguration2), new(circularParentConfiguration2), new(circularGrantConfiguration2)},
			err:            depends.ErrCircularDependency,
		},
	}

	for _, data := range testData {
		t.Run(data.title, func(t *testing.T) {
			_, err := depends.Resolve(data.configurations)
			assert.Equal(t, data.err, err)
		})
	}
}
