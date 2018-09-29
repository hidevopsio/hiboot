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
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/factory/depends/bar"
	"github.com/hidevopsio/hiboot/pkg/factory/depends/foo"
	"github.com/hidevopsio/hiboot/pkg/factory/depends/fake"
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
				{Kind: reflect.Ptr, Object: new(fooConfiguration)},
				{Kind: reflect.Ptr, Object: new(bar.Configuration)},
				{Kind: reflect.Ptr, Object: new(childConfiguration)},
				{Kind: reflect.Ptr, Object: new(fake.Configuration)},
				{Kind: reflect.Ptr, Object: new(parentConfiguration)},
				{Kind: reflect.Ptr, Object: new(grantConfiguration)},
				{Kind: reflect.Func, Object: foo.NewConfiguration},
				{Kind: reflect.Ptr, Object: new(barConfiguration)},
			},
			err: nil,
		},
		{
			title: "should sort dependencies",
			configurations: []*factory.MetaData{
				{Kind: reflect.Ptr, Object: new(fake.Configuration)},
				{Kind: reflect.Ptr, Object: new(fooConfiguration)},
				{Kind: reflect.Ptr, Object: new(bar.Configuration)},
				{Kind: reflect.Ptr, Object: new(childConfiguration)},
				{Kind: reflect.Ptr, Object: new(parentConfiguration)},
				{Kind: reflect.Ptr, Object: new(grantConfiguration)},
				{Kind: reflect.Func, Object: foo.NewConfiguration},
				{Kind: reflect.Ptr, Object: new(barConfiguration)},
			},
			err: nil,
		},
		//{
		//	title: "should report some of the dependencies are not found",
		//	configurations: []*factory.MetaData{
		//		{Kind: reflect.Ptr, Object: new(fooConfiguration)},
		//		{Kind: reflect.Ptr, Object: new(childConfiguration)},
		//		{Kind: reflect.Ptr, Object: new(parentConfiguration)},
		//		{Kind: reflect.Ptr, Object: new(grantConfiguration)},
		//		{Kind: reflect.Ptr, Object: new(barConfiguration)},
		//	},
		//	err: depends.ErrCircularDependency,
		//},
		{
			title: "should sort with constructor's dependencies",
			configurations: []*factory.MetaData{
				{Kind: reflect.Func, PkgName: "depends_test", Name: "barService", Object: newBarService},
				{Kind: reflect.Ptr, Object: new(Bar)},
				{Kind: reflect.Func, Object: newFooService},
				{Kind: reflect.Ptr, Object: new(Foo)},
				{Kind: reflect.Ptr, Object: new(Baz)},
				{Kind: reflect.Func, Object: newBazService},
			},
			err: nil,
		},
		{
			title: "should fail to sort with circular dependencies",
			configurations: []*factory.MetaData{
				{Kind: reflect.Ptr, Object: new(circularChildConfiguration)},
				{Kind: reflect.Ptr, Object: new(circularParentConfiguration)},
				{Kind: reflect.Ptr, Object: new(circularGrantConfiguration)},
			},
			err: depends.ErrCircularDependency,
		},
		{
			title: "should fail to sort with circular dependencies",
			configurations: []*factory.MetaData{
				{Kind: reflect.Ptr, Object: new(circularChildConfiguration2)},
				{Kind: reflect.Ptr, Object: new(circularParentConfiguration2)},
				{Kind: reflect.Ptr, Object: new(circularGrantConfiguration2)},
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
