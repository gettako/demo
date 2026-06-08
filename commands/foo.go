package commands

import (
	"flag"
	"github.com/takoterm/tako/internal/console/cli"
)

// FooCommand implements the cli.Command interface.
type FooCommand struct{}

func (c *FooCommand) Name() string {
	return "foo"
}

func (c *FooCommand) Description() string {
	return "Description for Foo command"
}

func (c *FooCommand) DefineFlags(fs *flag.FlagSet) {
	// fs.String("flag", "default", "description")
}

func (c *FooCommand) Execute(ctx *cli.Context, args []string) error {
	ctx.Info("Executing Foo command!")
	return nil
}
