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

func NewTestApplication(t *testing.T, cmd ...Command) TestApplication {
	a := new(testApplication)
	err := a.initialize(cmd...)
	assert.Equal(t, nil, err)
	a.build()
	return a
}

func (a *testApplication) RunTest(args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	a.root.SetOutput(buf)
	a.root.SetArgs(args)

	_, err = a.root.ExecuteC()

	return buf.String(), err
}
