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
	"testing"
	"github.com/hidevopsio/hiboot/pkg/log"
)

func init() {
	log.SetLevel(log.DebugLevel)
	AddCommand("root", new(demoCommand))
	AddCommand("root.demo", new(fooCommand))
	AddCommand("root.demo.foo", new(barCommand), new(bazCommand))
}

type demoCommand struct {
	BaseCommand
}

func (c *demoCommand) Init() {
	c.Use = "demo"
	c.Short = "demo command"
	c.Long = "Run demo command"
}

func (c *demoCommand) Handle(args []string) (err error) {
	log.Debug("here is demo command")
	return
}

type FooFlags struct {
	Name   string
	IntVal int
}

type fooCommand struct {
	BaseCommand
	//Short string `value:"foo sub command"`
	//Long  string `value:"Run foo sub command of command demo"`
}


func (c *fooCommand) Init() {
	c.Use = "foo"
	c.Short = "foo command"
	c.Long = "Run foo command"
}

func (c *fooCommand) Handle(args []string) (err error) {
	return nil
}

type barCommand struct {
	BaseCommand
	//Short string `value:"bar sub command"`
	//Long  string `value:"Run bar sub command of command foo"`
}


func (c *barCommand) Init() {
	c.Use = "bar"
	c.Short = "bar command"
	c.Long = "Run bar command"
}

func (c *barCommand) Handle(args []string) (err error) {
	return nil
}

type bazCommand struct {
	BaseCommand
	//Short string `value:"baz sub command"`
	//Long  string `value:"Run baz sub command of command foo"`
}


func (c *bazCommand) Init() {
	c.Use = "baz"
	c.Short = "baz command"
	c.Long = "Run baz command"
}

func (c *bazCommand) Handle(args []string) (err error) {
	return nil
}

// demo foo bar

func TestCliApplication(t *testing.T) {
	NewApplication().Run()
}
