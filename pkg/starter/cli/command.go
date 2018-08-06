package cli

import (
	"github.com/spf13/cobra"
	"errors"
)

type Command interface {
	EmbeddedCommand() *cobra.Command
	AddChild(child Command)
	Exec() error
	Name() string
	FullName() string
	SetName(name string)
	SetFullName(name string)
	Parent() Command
	SetParent(p Command)
	Handle(args []string) error
	Find(name string) (Command, error)
}


var CommandNotFoundError = errors.New("command not found")

type BaseCommand struct {
	cobra.Command
	name string
	fullname string
	parent Command
	children []Command
	childrenMap map[string]Command
}


func (c *BaseCommand) EmbeddedCommand() *cobra.Command  {
	return &c.Command
}

func (c *BaseCommand) Handle(args []string) error {
	return nil
}

func (c *BaseCommand) Exec() error {
	return c.Execute()
}

func (c *BaseCommand) AddChild(child Command) {
	childEmbeddedCommand := child.EmbeddedCommand()
	childEmbeddedCommand.RunE = func(cmd *cobra.Command, args []string) error {
		return child.Handle(args)
	}
	child.SetParent(c)
	c.children = append(c.children, child)
	c.AddCommand(child.EmbeddedCommand())
}

func (c *BaseCommand) Name() string {
	return c.name
}

func (c *BaseCommand) SetName(name string) {
	c.name = name
}

func (c *BaseCommand) FullName() string {
	return c.fullname
}

func (c *BaseCommand) SetFullName(name string) {
	c.fullname = name
}

func (c *BaseCommand) Parent() Command  {
	return c.parent
}

func (c *BaseCommand) SetParent(p Command)   {
	c.parent = p
}

func (c *BaseCommand) Find(name string) (Command, error)  {
	if c.name == name {
		return c, nil
	}

	for _, cmd := range c.children {
		if name == cmd.Name() {
			return cmd, nil
		}
	}
	return nil, CommandNotFoundError
}
