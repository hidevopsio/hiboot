package cmd

import (
	"github.com/hidevopsio/hiboot/pkg/starter/cli"
	"github.com/hidevopsio/hiboot/pkg/log"
)

// FirstCommand is the root command
type FirstCommand struct {
	// embedded cli.BaseCommand
	cli.BaseCommand

	// inject (add) secondCommand into FirstCommand
	Second *secondCommand `cmd:""`

	// inject flag
	Profile *string `flag:"shorthand=p,value=dev,usage=e.g. --profile=test"`
	Timeout *int `flag:"shorthand=t,value=1,usage=e.g. --timeout=2"`
}

func (c *FirstCommand) Init() {
	c.Use = "first"
	c.Short = "first command"
	c.Long = "Run first command"
	c.ValidArgs = []string{"baz"}
}

func (c *FirstCommand) Run(args []string) error {
	log.Infof("handle first command: profile=%v, timeout=%v", *c.Profile, *c.Timeout)
	return nil
}

