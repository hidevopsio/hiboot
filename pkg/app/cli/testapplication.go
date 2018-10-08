package cli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestApplication interface {
	Application
	RunTest(args ...string) (output string, err error)
}

type testApplication struct {
	application
}

// NewTestApplication is the test application constructor
func NewTestApplication(t *testing.T, cmd ...interface{}) TestApplication {
	a := new(testApplication)
	err := a.initialize(cmd...)
	assert.Equal(t, nil, err)
	err = a.build()
	assert.Equal(t, nil, err)
	return a
}

func (a *testApplication) RunTest(args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	if a.root != nil {
		a.root.SetOutput(buf)
		a.root.SetArgs(args)

		_, err = a.root.ExecuteC()

		return buf.String(), err
	}
	return
}
