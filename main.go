package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"

	. "github.com/cloudfoundry-incubator/stembuild/commandparser"
	"github.com/cloudfoundry-incubator/stembuild/version"
	. "github.com/google/subcommands"
)

//go:generate go run gen.go

func main() {
	var gf GlobalFlags
	packageCmd := PackageCmd{}
	packageCmd.GlobalFlags = &gf
	constructCmd := ConstructCmd{}
	constructCmd.GlobalFlags = &gf

	var commands = make([]Command, 0)

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.BoolVar(&gf.Debug, "debug", false, "Print lots of debugging information")
	fs.BoolVar(&gf.Color, "color", false, "Colorize debug output")
	fs.BoolVar(&gf.ShowVersion, "version", false, "Show Stembuild version")
	fs.BoolVar(&gf.ShowVersion, "v", false, "Stembuild version (shorthand)")

	commander := NewCommander(fs, path.Base(os.Args[0]))

	sh := NewStembuildHelp(commander, fs, &commands)
	commander.Register(sh, "")
	commands = append(commands, sh)

	commander.Register(&packageCmd, "")
	commander.Register(&constructCmd, "")

	commands = append(commands, &packageCmd)
	commands = append(commands, &constructCmd)

	// Override the default usage text of Google's Subcommand with our own
	fs.Usage = func() { sh.Explain(commander.Error) }

	_ = fs.Parse(os.Args[1:])
	if gf.ShowVersion {
		_, _ = fmt.Fprintf(os.Stdout, "%s version %s, Windows Stemcell Building Tool\n\n", path.Base(os.Args[0]), version.Version)
		os.Exit(0)
	}

	ctx := context.Background()
	os.Exit(int(commander.Execute(ctx)))
}
