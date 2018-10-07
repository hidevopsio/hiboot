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

package app_test

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestApp(t *testing.T) {
	type fakeProperties struct {
		Name string `default:"fake"`
	}
	type fakeConfiguration struct {
		app.Configuration
		Properties fakeProperties `mapstructure:"fake"`
	}
	t.Run("should add configuration", func(t *testing.T) {
		err := app.AutoConfiguration(new(fakeConfiguration))
		assert.Equal(t, nil, err)
	})

	//t.Run("should report duplication error", func(t *testing.T) {
	//	err := app.AutoConfiguration(new(fakeConfiguration))
	//	assert.Equal(t, app.ConfigurationNameIsTakenError, err)
	//})

	//t.Run("should not add invalid configuration", func(t *testing.T) {
	//	type fooConfiguration struct {
	//		Properties fakeProperties `mapstructure:"fake"`
	//	}
	//	err := app.AutoConfiguration(fooConfiguration{})
	//	assert.Equal(t, app.ErrInvalidObjectType, err)
	//})

	type configuration struct {
		app.PreConfiguration
		Properties fakeProperties `mapstructure:"fake"`
	}
	t.Run("should add configuration with pkg name", func(t *testing.T) {
		err := app.AutoConfiguration(new(configuration))
		assert.Equal(t, nil, err)
	})

	//t.Run("should add named configuration", func(t *testing.T) {
	//	err := app.AutoConfiguration("baz", new(configuration))
	//	assert.Equal(t, nil, err)
	//})

	t.Run("should not add invalid configuration", func(t *testing.T) {
		err := app.AutoConfiguration(nil)
		assert.Equal(t, app.ErrInvalidObjectType, err)
	})

	t.Run("should add configuration with pkg name", func(t *testing.T) {
		type configuration struct {
			app.PostConfiguration
			Properties fakeProperties `mapstructure:"fake"`
		}
		err := app.AutoConfiguration(new(configuration))
		assert.Equal(t, nil, err)
	})

	//t.Run("should not add invalid configuration which embedded unknown interface", func(t *testing.T) {
	//	type unknownInterface interface{}
	//	type configuration struct {
	//		unknownInterface
	//		Properties fakeProperties `mapstructure:"fake"`
	//	}
	//	err := app.AutoConfiguration(new(configuration))
	//	assert.Equal(t, app.InvalidObjectTypeError, err)
	//})

	//t.Run("should not add configuration with non point type", func(t *testing.T) {
	//	type configuration struct {
	//		app.Configuration
	//		Properties fakeProperties `mapstructure:"fake"`
	//	}
	//	err := app.AutoConfiguration(configuration{})
	//	assert.Equal(t, app.ErrInvalidObjectType, err)
	//})

	//t.Run("should not add invalid configuration that not embedded with app.Configuration", func(t *testing.T) {
	//	type invalidConfiguration struct {
	//		Properties fakeProperties `mapstructure:"fake"`
	//	}
	//	err := app.AutoConfiguration(new(invalidConfiguration))
	//	assert.Equal(t, app.ErrInvalidObjectType, err)
	//})

	t.Run("should not add invalid component", func(t *testing.T) {
		err := app.Component(nil)
		assert.Equal(t, app.ErrInvalidObjectType, err)
	})

	t.Run("should add new component", func(t *testing.T) {
		type fakeService interface{}
		type fakeServiceImpl struct{ fakeService }
		err := app.Component(new(fakeServiceImpl))
		assert.Equal(t, nil, err)
	})

	t.Run("should add new named component", func(t *testing.T) {
		type fakeService interface{}
		type fakeServiceImpl struct{ fakeService }
		err := app.Component("myService", new(fakeServiceImpl))
		assert.Equal(t, nil, err)
	})
}

func TestBaseApplication(t *testing.T) {
	ba := new(app.BaseApplication)

	ba.BeforeInitialization()

	err := ba.Initialize()
	assert.Equal(t, nil, err)

	sc := ba.SystemConfig()
	assert.NotEqual(t, nil, sc)

	// TODO: check concurrency issue during test
	ba.BuildConfigurations()

	ba.GetInstance("foo")

	cf := ba.ConfigurableFactory()
	assert.NotEqual(t, nil, cf)

	ba.AfterInitialization()

	ba.RegisterController(nil)

	ba.SetProperty(app.PropertyBannerDisabled, false).
		SetProperty(app.PropertyAppProfilesInclude, "foo")

	ba.AppendProfiles(ba)

	ba.PrintStartupMessages()

	ba.Use()

	ba.Run()

	ba.GetInstance("foo")

}
