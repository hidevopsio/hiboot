package app_test

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

	t.Run("should report duplication error", func(t *testing.T) {
		err := app.AutoConfiguration(new(fakeConfiguration))
		assert.Equal(t, app.ConfigurationNameIsTakenError, err)
	})

	t.Run("should not add invalid configuration", func(t *testing.T) {
		type fooConfiguration struct {
			Properties fakeProperties `mapstructure:"fake"`
		}
		err := app.AutoConfiguration(fooConfiguration{})
		assert.Equal(t, app.InvalidObjectTypeError, err)
	})

	type configuration struct {
		app.PreConfiguration
		Properties fakeProperties `mapstructure:"fake"`
	}
	t.Run("should add configuration with pkg name", func(t *testing.T) {
		err := app.AutoConfiguration(new(configuration))
		assert.Equal(t, nil, err)
	})

	t.Run("should add named configuration", func(t *testing.T) {
		err := app.AutoConfiguration("baz", new(configuration))
		assert.Equal(t, nil, err)
	})

	t.Run("should not add invalid configuration", func(t *testing.T) {
		err := app.AutoConfiguration(nil)
		assert.Equal(t, app.InvalidObjectTypeError, err)
	})

	t.Run("should add configuration with pkg name", func(t *testing.T) {
		type configuration struct {
			app.PostConfiguration
			Properties fakeProperties `mapstructure:"fake"`
		}
		err := app.AutoConfiguration(new(configuration))
		assert.Equal(t, nil, err)
	})

	t.Run("should not add invalid configuration which embedded unknown interface", func(t *testing.T) {
		type unknownInterface interface{}
		type configuration struct {
			unknownInterface
			Properties fakeProperties `mapstructure:"fake"`
		}
		err := app.AutoConfiguration(new(configuration))
		assert.Equal(t, app.InvalidObjectTypeError, err)
	})

	t.Run("should not add configuration with non point type", func(t *testing.T) {
		type configuration struct {
			app.Configuration
			Properties fakeProperties `mapstructure:"fake"`
		}
		err := app.AutoConfiguration(configuration{})
		assert.Equal(t, app.InvalidObjectTypeError, err)
	})

	t.Run("should not add invalid configuration that not embedded with app.Configuration", func(t *testing.T) {
		type invalidConfiguration struct {
			Properties fakeProperties `mapstructure:"fake"`
		}
		err := app.AutoConfiguration(new(invalidConfiguration))
		assert.Equal(t, app.InvalidObjectTypeError, err)
	})

	t.Run("should not add invalid component", func(t *testing.T) {
		err := app.Component(nil)
		assert.Equal(t, app.InvalidObjectTypeError, err)
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

	t.Run("should report component name collision", func(t *testing.T) {
		type fakeService interface{}
		type fakeServiceImpl struct{ fakeService }
		err := app.Component("myService", new(fakeServiceImpl))
		assert.Equal(t, app.ComponentNameIsTakenError, err)
	})
}

func TestBaseApplication(t *testing.T) {
	ba := new(app.BaseApplication)

	ba.BeforeInitialization()

	err := ba.Init()
	assert.Equal(t, nil, err)

	sc := ba.SystemConfig()
	assert.NotEqual(t, nil, sc)

	ba.BuildConfigurations()

	cf := ba.ConfigurableFactory()
	assert.NotEqual(t, nil, cf)

	ba.AfterInitialization()

	ba.RegisterController(nil)

	ba.SetProperty(app.PropertyBannerDisabled, false)

	ba.PrintStartupMessages()

	ba.Use()
}
