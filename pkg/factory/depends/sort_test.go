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
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/depends"
	"github.com/hidevopsio/hiboot/pkg/factory/depends/bar"
	"github.com/hidevopsio/hiboot/pkg/factory/depends/fake"
	"github.com/hidevopsio/hiboot/pkg/factory/depends/foo"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"github.com/magiconair/properties/assert"
	"reflect"
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

type circularChildConfiguration struct {
	app.Configuration `depends:"circularParentConfiguration"`
}

type circularParentConfiguration struct {
	app.Configuration `depends:"circularGrantConfiguration"`
}

type circularGrantConfiguration struct {
	app.Configuration `depends:"circularParentConfiguration"`
}

type circularChildConfiguration2 struct {
	app.Configuration `depends:"circularParentConfiguration2"`
}

type circularParentConfiguration2 struct {
	app.Configuration `depends:"circularGrantConfiguration2"`
}

type circularGrantConfiguration2 struct {
	app.Configuration `depends:"circularChildConfiguration2"`
}

type Foo struct {
	Name string
}

type Bar struct {
	Name string
}

type Baz struct {
	Name string
}

type fooService struct {
	foo *Foo
}

func newFooService(foo *Foo) *fooService {
	return &fooService{foo: foo}
}

type barService struct {
	bar        *Bar
	fooService *fooService
}

func newBarService(bar *Bar, fooService *fooService) *barService {
	return &barService{bar: bar, fooService: fooService}
}

type foobarService struct {
	fooService *fooService
}

func newFoobarService(fooService *fooService) *foobarService {
	return &foobarService{fooService: fooService}
}

type bazService struct {
	bar *Bar
	baz *Baz
}

func newBazService(bar *Bar, baz *Baz) *bazService {
	return &bazService{bar: bar, baz: baz}
}

func newData(object interface{}) *factory.MetaData {

	pkgName, name := reflector.GetPkgAndName(object)

	return &factory.MetaData{
		Kind: reflect.TypeOf(object).Kind(),
		// TODO: should check more conditions, like named instance,
		// var foobar *Foo and var foo *Foo should be supported
		//
		Name:     pkgName + "." + name,
		TypeName: name,
		Object:   object,
		PkgName:  pkgName,
	}
}

func TestSort(t *testing.T) {

	log.SetLevel(log.DebugLevel)

	testData := []struct {
		title          string
		configurations []*factory.MetaData
		err            error
	}{
		{
			title: "should sort dependencies",
			configurations: []*factory.MetaData{
				newData(new(fooConfiguration)),
				newData(new(bar.Configuration)),
				newData(new(fake.Configuration)),
				newData(new(parentConfiguration)),
				newData(new(grantConfiguration)),
				newData(new(childConfiguration)),
				newData(foo.NewConfiguration),
				newData(new(barConfiguration)),
			},
			err: nil,
		},
		{
			title: "should sort dependencies",
			configurations: []*factory.MetaData{
				newData(new(fake.Configuration)),
				newData(new(fooConfiguration)),
				newData(new(bar.Configuration)),
				newData(new(childConfiguration)),
				newData(new(grantConfiguration)),
				newData(new(parentConfiguration)),
				newData(foo.NewConfiguration),
				newData(new(barConfiguration)),
			},
			err: nil,
		},
		{
			title: "should report some of the dependencies are not found",
			configurations: []*factory.MetaData{
				newData(new(fooConfiguration)),
				newData(new(childConfiguration)),
				newData(new(grantConfiguration)),
				newData(new(parentConfiguration)),
				newData(new(barConfiguration)),
			},
			err: nil, // TODO: temp solution depends.ErrCircularDependency,
		},
		{
			title: "should sort with constructor's dependencies",
			configurations: []*factory.MetaData{
				newData(newBarService),
				newData(new(Bar)),
				newData(newFooService),
				newData(new(Foo)),
				newData(new(Baz)),
				newData(newBazService),
			},
			err: nil,
		},
		{
			title: "should fail to sort with circular dependencies",
			configurations: []*factory.MetaData{
				newData(new(circularChildConfiguration)),
				newData(new(circularParentConfiguration)),
				newData(new(circularGrantConfiguration)),
			},
			err: depends.ErrCircularDependency,
		},
		{
			title: "should fail to sort with circular dependencies",
			configurations: []*factory.MetaData{
				newData(new(Bar)),
				newData(new(circularChildConfiguration)),
				newData(new(circularParentConfiguration)),
				newData(new(circularGrantConfiguration)),
			},
			err: depends.ErrCircularDependency,
		},
		{
			title: "should fail to sort with circular dependencies",
			configurations: []*factory.MetaData{
				newData(new(circularChildConfiguration2)),
				newData(new(circularParentConfiguration2)),
				newData(new(circularGrantConfiguration2)),
			},
			err: depends.ErrCircularDependency,
		},
	}

	for _, data := range testData {
		t.Run(data.title, func(t *testing.T) {
			_, err := depends.Resolve(data.configurations)
			assert.Equal(t, err, data.err)
		})
	}
}
