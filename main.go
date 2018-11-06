package main

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	. "github.com/pivotal-cf-experimental/stembuild/commandparser"
	"os"
	"path"
)

func main() {
	var gf GlobalFlags

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.BoolVar(&gf.Debug, "debug", false, "Print lots of debugging information")
	fs.BoolVar(&gf.Color, "color", false, "Colorize debug output")

	commander := subcommands.NewCommander(fs, path.Base(os.Args[0]))

	commander.Register(commander.HelpCommand(), "")
	commander.Register(commander.FlagsCommand(), "")
	commander.Register(commander.CommandsCommand(), "")

	packageCmd := PackageCmd{}
	packageCmd.GlobalFlags = &gf
	commander.Register(&packageCmd, "")

	fs.Parse(os.Args[1:])
	ctx := context.Background()
	os.Exit(int(commander.Execute(ctx)))
}
