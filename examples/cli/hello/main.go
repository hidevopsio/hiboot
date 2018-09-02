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

// declare main package
package main

// import cli starter and fmt
import (
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/app/cli"
)

// define the command
type HelloCommand struct {
	// embedding cli.BaseCommand in each command
	cli.BaseCommand
	// inject (bind) flag to field 'To', so that it can be used on Run method, please note that the data type must be pointer
	To *string `flag:"name=to,shorthand=t,value=world,usage=e.g. --to=world or -t world"`
}

// Init constructor
func (c *HelloCommand) Init() {
	c.Use = "hello"
	c.Short = "hello command"
	c.Long = "run hello command for getting started"
	c.Example = `
hello -h : help
hello -t John : say hello to John
`
	c.BashCompletionFunction = `
`
}

// Run run the command
func (c *HelloCommand) Run(args []string) error {
	fmt.Printf("Hello, %v\n", *c.To)
	return nil
}

// main function
func main() {
	// create new cli application and run it
	cli.NewApplication(new(HelloCommand)).
		SetProperty("banner.disabled", true).
		Run()
}
