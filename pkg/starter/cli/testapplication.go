package cli

import (
	"bytes"
	"path/filepath"
	"os"
	"runtime"
	"strings"
)

type TestApplication interface {
	Application
	RunTest(args ...string) (output string, err error)
}

type testApplication struct {
	application
}

func NewTestApplication(cmd ...Command) TestApplication {
	basename := filepath.Base(os.Args[0])
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}
	var ta TestApplication
	ta = new(testApplication)
	root := new(rootCommand)
	root.SetName("root")
	if cmd != nil && len(cmd) > 0 {
		root.Add(cmd...)
	}
	ta.SetRoot(root)
	ta.Root().EmbeddedCommand().Use = basename
	ta.Init()
	return ta
}

func (a *testApplication) RunTest(args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	a.root.SetOutput(buf)
	a.root.SetArgs(args)

	_, err = a.root.ExecuteC()

	return buf.String(), err
}


