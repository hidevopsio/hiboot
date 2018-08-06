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

package cli

import (
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"reflect"
	"os"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"github.com/hidevopsio/hiboot/pkg/log"
	"strings"
	"sync"
	"runtime"
	"path/filepath"
	"github.com/hidevopsio/hiboot/pkg/utils/sort"
)

type Application interface {
	Run()
	Initialize() error
	Root() Command
	SetRoot(root Command)
	FindCommand(name string) Command
}

type application struct {
	rootCommand Command
}

type CommandNameValue struct {
	Name string
	Command interface{}
}

var (
	commandContainer  map[string][]Command
	commandNames []string
	app *application
)

func init() {
	commandContainer = make(map[string][]Command)
	commandNames = make([]string, 0)
}


/*

root(demo)
demo(foo)
foo(bar, baz)

demo - foo - bar
           - baz

foo.parentName = root
bar.parentName = foo
baz.parentName = foo

*/

func GetCommandName(cmd interface{}) string {
	name, err := reflector.GetName(cmd)
	if err == nil {
		name = strings.Replace(name, "Command", "", -1)
		name = strings.ToLower(name)
	}
	return name
}

func AddCommand(parentName string, commands ...Command) {
	commandNames = append(commandNames, parentName)
	for _, command := range commands {
		commandContainer[parentName] = append(commandContainer[parentName], command)
	}
}

var do sync.Once

func GetApplication() Application {
	do.Do(func() {
		app = new(application)
		app.SetRoot(new(rootCommand))
	})
	return app
}

func NewApplication() Application {
	basename := filepath.Base(os.Args[0])
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	// TODO: read config file, replace basename if user specified app.name

	a := GetApplication()
	a.Root().EmbeddedCommand().Use = basename
	a.Initialize()
	return a
}

func (a *application) Initialize() error  {
	// parse commands
	parentContainer := make(map[string]Command)
	sort.SortByLen(commandNames)
	fullname := "root"
	parentContainer[fullname] = a.rootCommand
	for _, cmdName := range commandNames {
		commands := commandContainer[cmdName]
		parent := parentContainer[cmdName]
		for _, command := range commands {
			inject.IntoObject(reflect.ValueOf(command))
			command.SetName(GetCommandName(command))
			fullname := cmdName + "." + command.Name()
			parentContainer[fullname] = command
			command.SetFullName(fullname)
			parent.AddChild(command)
		}
	}
	return nil
}

func (a *application) SetRoot(root Command) {
	a.rootCommand = root
}

func (a *application) Root() Command {
	return a.rootCommand
}

func (a *application) FindCommand(name string) Command {
	names := strings.SplitN(name, ".", -1)
	cmd := a.rootCommand
	foundCmd := cmd
	var err error
	for i, n := range names {
		if n == cmd.Name() {
			cmd, err = cmd.Find(names[i + 1])
			if err == CommandNotFoundError {
				break
			}
			foundCmd = cmd
		}
	}
	return foundCmd
}


func (a *application) Run() {

	log.Debug(commandContainer)

	if a.rootCommand != nil{
		if err := a.rootCommand.Exec(); err != nil {
			os.Exit(1)
		}
	}
}



