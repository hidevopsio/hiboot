package cmd

import (
	"github.com/hidevopsio/hiboot/pkg/starter/cli"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type firstCommand struct {
	cli.BaseCommand
}

func init() {
	cli.AddCommand("root", new(firstCommand))
}

func (c *firstCommand) Init() {
	c.Use = "first"
	c.Short = "first command"
	c.Long = "Run first command"
}

func (c *firstCommand) Handle(args []string) error {
	log.Debug("handle first command")
	return nil
}

