package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cloudfoundry-incubator/stembuild/assets"
	. "github.com/cloudfoundry-incubator/stembuild/commandparser"
	"github.com/cloudfoundry-incubator/stembuild/construct/factory"
	"github.com/cloudfoundry-incubator/stembuild/version"
	. "github.com/google/subcommands"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	data, err := assets.Asset("StemcellAutomation.zip")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "StemcellAutomation not found")
		os.Exit(1)
	}
	s := "./StemcellAutomation.zip"
	err = ioutil.WriteFile(s, data, 0644)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Unable to write StemcellAutomation.zip")
		os.Exit(1)
	}

	var gf GlobalFlags
	packageCmd := PackageCmd{}
	packageCmd.GlobalFlags = &gf
	constructCmd := NewConstructCmd(&vmconstruct_factory.VMConstructFactory{}, &ConstructValidator{}, &ConstructCmdMessenger{OutputChannel: os.Stderr})
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
		_ = os.Remove(s)
		os.Exit(0)
	}

	ctx := context.Background()
	i := int(commander.Execute(ctx))
	_ = os.Remove(s)
	os.Exit(i)
}
