package cli

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/cobra"
	"errors"
	"io"
)

type Command interface {
	EmbeddedCommand() *cobra.Command
	Add(commands ...Command) Command
	HasChild() bool
	Children() []Command
	Exec() error
	Name() string
	FullName() string
	SetName(name string) Command
	SetFullName(name string) Command
	Parent() Command
	SetParent(p Command) Command
	Run(args []string) error
	Find(name string) (Command, error)
	SetOutput(output io.Writer)
	SetArgs(a []string)
	ExecuteC() (cmd *cobra.Command, err error)
	PersistentFlags() *flag.FlagSet
}

var CommandNotFoundError = errors.New("command not found")

type BaseCommand struct {
	cobra.Command
	name string
	fullname string
	parent Command
	children []Command
}

func (c *BaseCommand) EmbeddedCommand() *cobra.Command  {
	return &c.Command
}

func (c *BaseCommand) Run(args []string) error {
	return nil
}

func (c *BaseCommand) Exec() error {
	return c.Execute()
}

func (c *BaseCommand) HasChild() bool {
	return len(c.children) > 0
}

func (c *BaseCommand) Children() []Command {
	return c.children
}

func (c *BaseCommand) addChild(child Command) {
	if child.Name() == "" {
		name := parseName(child)
		child.SetName(name)
	}
	childEmbeddedCommand := child.EmbeddedCommand()
	childEmbeddedCommand.RunE = func(cmd *cobra.Command, args []string) error {
		return child.Run(args)
	}
	child.SetParent(c)
	c.children = append(c.children, child)
	c.AddCommand(child.EmbeddedCommand())
}

func (c *BaseCommand) Add(commands ...Command) Command {
	for _, command := range commands {
		c.addChild(command)
	}
	return c
}

func (c *BaseCommand) Name() string {
	return c.name
}

func (c *BaseCommand) SetName(name string) Command {
	c.name = name
	return c
}

func (c *BaseCommand) FullName() string {
	if c.fullname == "" {
		c.fullname = c.name
	}
	return c.fullname
}

func (c *BaseCommand) SetFullName(name string) Command {
	c.fullname = name
	return c
}

func (c *BaseCommand) Parent() Command  {
	return c.parent
}

func (c *BaseCommand) SetParent(p Command) Command {
	c.parent = p
	return c
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
