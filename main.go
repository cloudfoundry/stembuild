package main

import (
	"context"
	"flag"
	. "github.com/google/subcommands"
	. "github.com/pivotal-cf-experimental/stembuild/commandparser"
	"os"
	"path"
)

func main() {
	var gf GlobalFlags
	packageCmd := PackageCmd{}
	packageCmd.GlobalFlags = &gf

	var commands = make([]Command, 0)

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.BoolVar(&gf.Debug, "debug", false, "Print lots of debugging information")
	fs.BoolVar(&gf.Color, "color", false, "Colorize debug output")

	commander := NewCommander(fs, path.Base(os.Args[0]))

	sh := NewStembuildHelp(commander, fs, &commands)
	commander.Register(sh, "")
	commands = append(commands, sh)

	commander.Register(&packageCmd, "")
	commands = append(commands, &packageCmd)

	// Override the default usage text of Google's Subcommand with our own
	fs.Usage = func() { sh.Explain() }

	_ = fs.Parse(os.Args[1:])
	ctx := context.Background()
	os.Exit(int(commander.Execute(ctx)))
}
