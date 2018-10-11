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

package instantiate_test

import (
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"github.com/deckarep/golang-set"
	"fmt"
)

type FooBar struct {
	Name string
}

type fooBarService struct {
	Name   string
	fooBar *FooBar
}

type BarService interface {
	Baz() string
}

type BarServiceImpl struct {
	BarService
}

func (s *BarServiceImpl) Baz() string {
	return "baz"
}

func newFooBarService(fooBar *FooBar) *fooBarService {
	return &fooBarService{
		fooBar: fooBar,
	}
}

type helloService struct {
}

type HelloWorld string

func (s *helloService) HelloWorld() HelloWorld {
	return HelloWorld("Hello world")
}

func TestInstantiateFactory(t *testing.T) {
	type foo struct{ Name string }
	f := new(foo)

	appFactory := new(instantiate.InstantiateFactory)
	testName := "foobar"
	t.Run("should failed to set/get instance when factory is not initialized", func(t *testing.T) {
		inst := appFactory.GetInstance("not-exist-instance")
		assert.Equal(t, nil, inst)

		err := appFactory.SetInstance("foo", nil)
		assert.Equal(t, instantiate.ErrNotInitialized, err)

		item := appFactory.Items()
		assert.Equal(t, 0, len(item))
	})

	hello := new(helloService)
	helloTyp := reflect.TypeOf(hello)
	numOfMethod := helloTyp.NumMethod()
	//log.Debug("methods: ", numOfMethod)
	testComponents := make([]*factory.MetaData, 0)
	for mi := 0; mi < numOfMethod; mi++ {
		method := helloTyp.Method(mi)
		// append inst to f.components
		testComponents = append(testComponents, factory.NewMetaData(hello, method))
	}

	testComponents = append(testComponents, factory.NewMetaData(f),
		factory.NewMetaData(&FooBar{Name: testName}),
		factory.NewMetaData(&BarServiceImpl{}),
		factory.NewMetaData(newFooBarService))

	ic := cmap.New()
	t.Run("should initialize factory", func(t *testing.T) {
		appFactory.Initialize(ic, testComponents)
		assert.Equal(t, true, appFactory.Initialized())
	})

	t.Run("should build components", func(t *testing.T) {
		appFactory.BuildComponents()
	})

	t.Run("should get built instance", func(t *testing.T) {
		inst := appFactory.GetInstance("instantiate_test.helloWorld")
		assert.Equal(t, HelloWorld("Hello world"), inst)
	})

	t.Run("should set and get instance from factory", func(t *testing.T) {
		appFactory.SetInstance(f)
		inst := appFactory.GetInstance(foo{})
		assert.Equal(t, f, inst)
	})

	t.Run("should failed to get instance that does not exist", func(t *testing.T) {
		inst := appFactory.GetInstance("not-exist-instance")
		assert.Equal(t, nil, inst)
	})

	t.Run("should failed to get instances that does not exist", func(t *testing.T) {
		inst := appFactory.GetInstances("not-exist-instances")
		assert.Equal(t, 0, len(inst))
	})

	t.Run("should failed to set instance that it already exists in test mode", func(t *testing.T) {
		nf := new(foo)
		err := appFactory.SetInstance("foo", nf)
		assert.NotEqual(t, nil, err)
	})

	t.Run("should get factory items", func(t *testing.T) {
		items := appFactory.Items()
		assert.NotEqual(t, 0, len(items))
	})

	t.Run("should check valid object", func(t *testing.T) {
		assert.Equal(t, true, appFactory.IsValidObjectType(f))
	})

	t.Run("should check invalid object", func(t *testing.T) {
		assert.Equal(t, false, appFactory.IsValidObjectType(1))
	})

	t.Run("should append component", func(t *testing.T) {
		appFactory.AppendComponent(&struct {
			Name string
		}{Name: "foo"})
		assert.Equal(t, false, appFactory.IsValidObjectType(1))
	})
}

func TestMapSet(t *testing.T) {
	requiredClasses := mapset.NewSet()
	requiredClasses.Add("Cooking")
	requiredClasses.Add("English")
	requiredClasses.Add("Math")
	requiredClasses.Add("Biology")

	scienceSlice := []interface{}{"Biology", "Chemistry"}
	scienceClasses := mapset.NewSetFromSlice(scienceSlice)

	electiveClasses := mapset.NewSet()
	electiveClasses.Add("Welding")
	electiveClasses.Add("Music")
	electiveClasses.Add("Automotive")

	bonusClasses := mapset.NewSet()
	bonusClasses.Add("Go Programming")
	bonusClasses.Add("Python Programming")

	//Show me all the available classes I can take
	allClasses := requiredClasses.Union(scienceClasses).Union(electiveClasses).Union(bonusClasses)
	fmt.Println(allClasses)

	//Is cooking considered a science class?
	fmt.Println(scienceClasses.Contains("Cooking")) //false

	//Show me all classes that are not science classes, since I hate science.
	fmt.Println(allClasses.Difference(scienceClasses)) //Set{Music, Automotive, Go Programming, Python Programming, Cooking, English, Math, Welding}

	//Which science classes are also required classes?
	fmt.Println(scienceClasses.Intersect(requiredClasses)) //Set{Biology}

	//How many bonus classes do you offer?
	fmt.Println(bonusClasses.Cardinality()) //2

	//Do you have the following classes? Welding, Automotive and English?
	fmt.Println(allClasses.IsSuperset(mapset.NewSetFromSlice([]interface{}{"Welding", "Automotive", "English"}))) //true
}