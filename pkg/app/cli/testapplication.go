package cli

import (
	"bytes"
)

type TestApplication interface {
	Application
	RunTest(args ...string) (output string, err error)
}

type testApplication struct {
	application
}

func NewTestApplication(cmd ...Command) TestApplication {
	ta := new(testApplication)
	ta.Init(cmd...)
	return ta
}

func (a *testApplication) RunTest(args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	a.root.SetOutput(buf)
	a.root.SetArgs(args)

	_, err = a.root.ExecuteC()

	return buf.String(), err
}


