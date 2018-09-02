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

// Package cli provides quick start for command line application
package cli

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/utils/gotest"
	"github.com/hidevopsio/hiboot/pkg/utils/sort"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"github.com/hidevopsio/hiboot/pkg/log"
)

// Application cli application interface
type Application interface {
	app.Application
	Root() Command
	SetRoot(root Command)
}

type application struct {
	app.BaseApplication
	root Command
}

// CommandNameValue
type CommandNameValue struct {
	Name    string
	Command interface{}
}

var (
	commandContainer map[string][]Command
	commandNames     []string
	cliApp           *application
	once             sync.Once
)

func init() {
	commandContainer = make(map[string][]Command)
	commandNames = make([]string, 0)
}

// HideBanner hide banner display on application start up
func HideBanner() {
	app.HideBanner()
}

// AddCommand add new command
func AddCommand(parentPath string, commands ...Command) {
	// de-duplication
	if commandContainer[parentPath] == nil {
		commandNames = append(commandNames, parentPath)
	}
	for _, command := range commands {
		commandContainer[parentPath] = append(commandContainer[parentPath], command)
	}
}

// NewApplication create new cli application
func NewApplication(cmd ...Command)  Application {
	a := new(application)
	if a.initialize(cmd...) != nil {
		log.Fatal("cli application is not initialized")
		os.Exit(1)
	}
	return a
}

func (a *application) injectCommand(cmd Command) {
	fullname := "root"
	if cmd != nil {
		fullname = cmd.FullName()
	}
	for _, child := range cmd.Children() {
		inject.IntoObjectValue(reflect.ValueOf(child))
		child.SetFullName(fullname + "." + child.GetName())
		a.injectCommand(child)
	}
}

func (a *application) initialize(cmd ...Command) (err error) {
	err = a.Init()
	if err == nil {
		var root Command
		root = new(rootCommand)
		numOfCmd := len(cmd)
		if cmd != nil && numOfCmd > 0 {
			if numOfCmd == 1 {
				root = cmd[0]
			} else {
				root.Add(cmd...)
			}
		}
		root.SetName("root")
		a.SetRoot(root)
	}
	return
}

// Init initialize cli application
func (a *application) build() error {
	basename := filepath.Base(os.Args[0])
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	var root = a.Root()
	inject.IntoObject(root)
	Register(root)
	a.SetRoot(root)
	if !gotest.IsRunning() {
		a.Root().EmbeddedCommand().Use = basename
	}

	if a.root != nil && a.root.HasChild() {
		a.injectCommand(a.root)
	} else {
		// parse commands
		parentContainer := make(map[string]Command)
		fullname := "root"
		sort.SortByLen(commandNames)
		parentContainer[fullname] = a.root
		for _, cmdName := range commandNames {
			commands := commandContainer[cmdName]
			parent := parentContainer[cmdName]
			if parent == nil {
				parent = a.root
			}
			for _, command := range commands {
				inject.IntoObjectValue(reflect.ValueOf(command))
				parent.Add(command)
				fullname := cmdName + "." + command.GetName()
				parentContainer[fullname] = command
				command.SetFullName(fullname)
			}
		}
	}
	return nil
}

// SetRoot set root command
func (a *application) SetRoot(root Command) {
	a.root = root
}

// Root get the root command
func (a *application) Root() Command {
	return a.root
}

// Run run the cli application
func (a *application) Run() {
	a.build()
	//log.Debug(commandContainer)
	if a.root != nil {
		if err := a.root.Exec(); err != nil {
			os.Exit(1)
		}
	}
}
