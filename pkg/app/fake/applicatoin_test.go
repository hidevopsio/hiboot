package fake

import "testing"

func TestApplicationContext(t *testing.T) {
	ac := new(ApplicationContext)
	ac.RegisterController(nil)
	ac.Use()
	ac.GetProperty("foo")
}