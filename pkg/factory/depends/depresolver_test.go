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
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type Hello string
type HelloWorld string
type HelloHiboot string

type helloConfiguration struct {
	app.Configuration
}

func (c *helloConfiguration) Hello() Hello {
	return Hello("Hello World")
}

func (c *helloConfiguration) HelloWorld(h Hello) HelloWorld {
	return HelloWorld(h + "World")
}

func (c *helloConfiguration) HelloHibootWorld(h Hello) HelloWorld {
	return HelloWorld(h + "Hiboot world")
}

func (c *helloConfiguration) HelloHiboot(h Hello) HelloHiboot {
	return HelloHiboot(h + "Hello Hiboot")
}

type helloService struct {
	HelloWorld  HelloWorld  `inject:""`
	HibootWorld HelloWorld  `inject:""`
	HelloHiboot HelloHiboot `inject:""`
}

type fooConfiguration struct {
	app.Configuration
}

func newFooConfiguration() *fooConfiguration {
	return &fooConfiguration{}
}

type barConfiguration struct {
	app.Configuration
}

func newBarConfiguration() *barConfiguration {
	return &barConfiguration{}
}

type childConfiguration struct {
	app.Configuration

	parent *parentConfiguration
}

func newChildConfiguration(parent *parentConfiguration) *childConfiguration {
	return &childConfiguration{parent: parent}
}

type parentConfiguration struct {
	app.Configuration

	grant *grantConfiguration
}

func newParentConfiguration(grant *grantConfiguration) *parentConfiguration {
	return &parentConfiguration{
		grant: grant,
	}
}

type grantConfiguration struct {
	app.Configuration

	fakeCfg *fake.Configuration
}

func newGrantConfiguration(fakeCfg *fake.Configuration) *grantConfiguration {
	return &grantConfiguration{fakeCfg: fakeCfg}
}

type circularChildConfiguration struct {
	app.Configuration

	circular *circularChildConfiguration
}

func newCircularChildConfiguration(circular *circularChildConfiguration) *circularChildConfiguration {
	return &circularChildConfiguration{circular: circular}
}

type circularParentConfiguration struct {
	app.Configuration
	circular *circularParentConfiguration
}

func newCircularParentConfiguration(circular *circularParentConfiguration) *circularParentConfiguration {
	return &circularParentConfiguration{circular: circular}
}

type circularGrantConfiguration struct {
	app.Configuration
	circular *circularGrantConfiguration
}

func newCircularGrantConfiguration(circular *circularGrantConfiguration) *circularGrantConfiguration {
	return &circularGrantConfiguration{circular: circular}
}

type circularChildConfiguration2 struct {
	app.Configuration
	circular *circularChildConfiguration2
}

func newCircularChildConfiguration2(circular *circularChildConfiguration2) *circularChildConfiguration2 {
	return &circularChildConfiguration2{circular: circular}
}

type circularParentConfiguration2 struct {
	app.Configuration
	circular *circularParentConfiguration2
}

func newCircularParentConfiguration2(circular *circularParentConfiguration2) *circularParentConfiguration2 {
	return &circularParentConfiguration2{circular: circular}
}

type circularGrantConfiguration2 struct {
	app.Configuration
	circular *circularGrantConfiguration2
}

func newCircularGrantConfiguration2(circular *circularGrantConfiguration2) *circularGrantConfiguration2 {
	return &circularGrantConfiguration2{circular: circular}
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

	helloConfig := new(helloConfiguration)
	helloTyp := reflect.TypeOf(helloConfig)
	numOfMethod := helloTyp.NumMethod()
	//log.Debug("methods: ", numOfMethod)
	methodTestData := make([]*factory.MetaData, 0)
	for mi := 0; mi < numOfMethod; mi++ {
		method := helloTyp.Method(mi)
		// append inst to f.components
		methodTestData = append(methodTestData, factory.NewMetaData(helloConfig, method))
	}
	methodTestData = append(methodTestData, factory.NewMetaData(helloConfig, new(helloService)))

	testData := []struct {
		title          string
		configurations []*factory.MetaData
		err            error
	}{
		{
			title:          "should sort with method's dependencies",
			configurations: methodTestData,
		},
		{
			title: "should sort dependencies",
			configurations: []*factory.MetaData{
				factory.NewMetaData(bar.NewConfiguration),
				factory.NewMetaData(foo.NewConfiguration),
				factory.NewMetaData(fake.NewConfiguration),
				factory.NewMetaData(newParentConfiguration),
				factory.NewMetaData(newGrantConfiguration),
				factory.NewMetaData(newChildConfiguration),
				factory.NewMetaData(foo.NewConfiguration),
			},
			err: nil,
		},
		{
			title: "should sort dependencies",
			configurations: []*factory.MetaData{
				factory.NewMetaData(fake.NewConfiguration),
				factory.NewMetaData(newFooConfiguration),
				factory.NewMetaData(bar.NewConfiguration),
				factory.NewMetaData(newChildConfiguration),
				factory.NewMetaData(newGrantConfiguration),
				factory.NewMetaData(newParentConfiguration),
				factory.NewMetaData(foo.NewConfiguration),
				factory.NewMetaData(bar.NewConfiguration),
			},
			err: nil,
		},
		{
			title: "should report some of the dependencies are not found",
			configurations: []*factory.MetaData{
				factory.NewMetaData(newFooConfiguration),
				factory.NewMetaData(newChildConfiguration),
				factory.NewMetaData(newGrantConfiguration),
				factory.NewMetaData(newParentConfiguration),
				factory.NewMetaData(newBarConfiguration),
			},
			err: depends.ErrCircularDependency,
		},
		{
			title: "should sort with constructor's dependencies",
			configurations: []*factory.MetaData{
				factory.NewMetaData(newBarService),
				factory.NewMetaData(new(Bar)),
				factory.NewMetaData(newFooService),
				factory.NewMetaData(new(Foo)),
				factory.NewMetaData(new(Baz)),
				factory.NewMetaData(newBazService),
			},
			err: nil,
		},
		{
			title: "should fail to sort with circular dependencies 1",
			configurations: []*factory.MetaData{
				factory.NewMetaData(newCircularChildConfiguration),
				factory.NewMetaData(newCircularParentConfiguration),
				factory.NewMetaData(newCircularGrantConfiguration),
			},
			err: depends.ErrCircularDependency,
		},
		{
			title: "should fail to sort with circular dependencies 2",
			configurations: []*factory.MetaData{
				factory.NewMetaData(new(Bar)),
				factory.NewMetaData(newCircularChildConfiguration),
				factory.NewMetaData(newCircularParentConfiguration),
				factory.NewMetaData(newCircularGrantConfiguration),
			},
			err: depends.ErrCircularDependency,
		},
		{
			title: "should fail to sort with circular dependencies 3",
			configurations: []*factory.MetaData{
				factory.NewMetaData(newCircularChildConfiguration2),
				factory.NewMetaData(newCircularParentConfiguration2),
				factory.NewMetaData(newCircularGrantConfiguration2),
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
