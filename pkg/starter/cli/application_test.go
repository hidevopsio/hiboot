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
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
	AddCommand("root", new(demoCommand))
	AddCommand("root.demo", new(fooCommand))
	AddCommand("root.demo.foo", new(barCommand), new(bazCommand))
}

type demoCommand struct {
	BaseCommand
	Profile *string `flag:"shorthand=p,value=dev,usage=e.g. --profile=test"`
	IntVal *int `flag:"name=integer,shorthand=i,value=0,usage=e.g. --integer=1"`
}

func (c *demoCommand) Init() {
	c.Use = "demo"
	c.Short = "demo command"
	c.Long = "Run demo command"
}

func (c *demoCommand) Run(args []string) (err error) {
	log.Debugf("on demo command - profile: %v, intVal: %v", *c.Profile, *c.IntVal)
	return
}

type fooCommand struct {
	BaseCommand
}


func (c *fooCommand) Init() {
	c.Use = "foo"
	c.Short = "foo command"
	c.Long = "Run foo command"
}

func (c *fooCommand) Run(args []string) (err error) {
	log.Debug("on foo command")
	return nil
}

type barCommand struct {
	BaseCommand
}

func (c *barCommand) Init() {
	c.Use = "bar"
	c.Short = "bar command"
	c.Long = "Run bar command"
}

func (c *barCommand) Run(args []string) (err error) {
	log.Debug("on bar command")
	return nil
}

type bazCommand struct {
	BaseCommand
}

func (c *bazCommand) Init() {
	c.Use = "baz"
	c.Short = "baz command"
	c.Long = "Run baz command"
}

func (c *bazCommand) Run(args []string) (err error) {
	log.Debug("on baz command")
	return nil
}

// demo foo bar
func TestCliApplication(t *testing.T) {

	demoCmd := new(demoCommand)
	fooCmd := new(fooCommand)
	fooCmd.Add(new(barCommand), new(bazCommand))
	demoCmd.Add(fooCmd)

	testApp := NewTestApplication(demoCmd)

	t.Run("should run root command", func(t *testing.T) {
		_, err := testApp.RunTest("-p", "test", "-i", "2")
		assert.Equal(t, nil, err)
	})

	t.Run("should run foo command", func(t *testing.T) {
		_, err := testApp.RunTest("foo")
		assert.Equal(t, nil, err)
	})

	t.Run("should run bar command", func(t *testing.T) {
		_, err := testApp.RunTest("foo", "bar")
		assert.Equal(t, nil, err)
	})

	t.Run("should run baz command", func(t *testing.T) {
		_, err := testApp.RunTest("foo", "baz")
		assert.Equal(t, nil, err)
	})
}

func TestNewApplication(t *testing.T) {
	go NewApplication().Run()
}

type A struct{
	Name string
}

func (a *A) Run(x string, y int) {
	log.Debugf("name: %v, x: %v, y: %v", a.Name, x, y)
}

type B struct{
	A
}

func (b *B) Run(x string)  {
	log.Debugf("x: %v", x)
	b.A.Run(x, 123)
}

func TestAB(t *testing.T) {
	b := new(B)
	b.Name = "bb"
	b.Run("hello")
}