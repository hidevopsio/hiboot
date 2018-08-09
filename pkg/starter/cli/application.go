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
	"reflect"
	"os"
	"github.com/hidevopsio/hiboot/pkg/inject"
	"strings"
	"sync"
	"runtime"
	"path/filepath"
	"github.com/hidevopsio/hiboot/pkg/utils/sort"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
)

type Application interface {
	Run()
	Init() error
	Root() Command
	SetRoot(root Command)
}

type application struct {
	root Command
}

type CommandNameValue struct {
	Name string
	Command interface{}
}

var (
	commandContainer  map[string][]Command
	commandNames []string
	app *application
	once sync.Once
)

func init() {
	commandContainer = make(map[string][]Command)
	commandNames = make([]string, 0)
}

func AddCommand(parentName string, commands ...Command) {
	commandNames = append(commandNames, parentName)
	for _, command := range commands {
		commandContainer[parentName] = append(commandContainer[parentName], command)
	}
}

func GetApplication() Application {
	once.Do(func() {
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
	a.Init()
	return a
}


func parseName(cmd Command) string {
	name, err := reflector.GetName(cmd)
	if err == nil {
		name = strings.Replace(name, "Command", "", -1)
		name = strings.ToLower(name)
	}
	return name
}

func (a *application) injectCommand(cmd Command)  {
	fullname := "root"
	if cmd != nil {
		fullname = cmd.FullName()
	}
	for _, child := range cmd.Children() {
		inject.IntoObject(reflect.ValueOf(child))
		child.SetFullName(fullname + "." + child.Name())
		a.injectCommand(child)
	}
}

func (a *application) Init() error  {
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
				inject.IntoObject(reflect.ValueOf(command))
				parent.Add(command)
				fullname := cmdName + "." + command.Name()
				parentContainer[fullname] = command
				command.SetFullName(fullname)
			}
		}
	}
	return nil
}

func (a *application) SetRoot(root Command) {
	a.root = root
}

func (a *application) Root() Command {
	return a.root
}

func (a *application) Run() {

	//log.Debug(commandContainer)

	if a.root != nil{
		if err := a.root.Exec(); err != nil {
			os.Exit(1)
		}
	}
}



