package cmd

import (
	"github.com/hidevopsio/hiboot/pkg/starter/cli"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type firstCommand struct {
	cli.BaseCommand
	Profile *string `flag:"shorthand=p,value=dev,usage=e.g. --profile=test"`
	Timeout *int `flag:"shorthand=t,value=1,usage=e.g. --timeout=2"`
}

func init() {
	cli.AddCommand("root", new(firstCommand))
}

func (c *firstCommand) Init() {
	c.Use = "first"
	c.Short = "first command"
	c.Long = "Run first command"
}

func (c *firstCommand) Run(args []string) error {
	log.Infof("handle first command: profile=%v, timeout=%v", *c.Profile, *c.Timeout)
	return nil
}

