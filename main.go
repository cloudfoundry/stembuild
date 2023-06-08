package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/cloudfoundry/stembuild/assets"
	. "github.com/cloudfoundry/stembuild/commandparser"
	vmconstruct_factory "github.com/cloudfoundry/stembuild/construct/factory"
	vcenter_client_factory "github.com/cloudfoundry/stembuild/iaas_cli/iaas_clients/factory"
	packager_factory "github.com/cloudfoundry/stembuild/package_stemcell/factory"
	"github.com/cloudfoundry/stembuild/version"
	. "github.com/google/subcommands"
)

func main() {
	envs := os.Environ()
	for _, env := range envs {
		env_name := strings.Split(env, "=")[0]
		if strings.HasPrefix(env_name, "GOVC_") || strings.HasPrefix(env_name, "GOVMOMI_") {
			fmt.Fprintf(os.Stderr, "Warning: The following environment variable is set and might override flags provided to stembuild: %s\n", env_name)
		}
	}
	data, err := assets.Asset("StemcellAutomation.zip")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "StemcellAutomation not found")
		os.Exit(1)
	}
	s := "./StemcellAutomation.zip"
	err = os.WriteFile(s, data, 0644)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Unable to write StemcellAutomation.zip")
		os.Exit(1)
	}

	var gf GlobalFlags
	packageCmd := NewPackageCommand(version.NewVersionGetter(), &packager_factory.PackagerFactory{}, &PackageMessenger{Output: os.Stderr})
	packageCmd.GlobalFlags = &gf
	constructCmd := NewConstructCmd(context.Background(), &vmconstruct_factory.VMConstructFactory{}, &vcenter_client_factory.ManagerFactory{}, &ConstructValidator{}, &ConstructCmdMessenger{OutputChannel: os.Stderr})
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

	commander.Register(packageCmd, "")
	commander.Register(constructCmd, "")

	commands = append(commands, packageCmd)
	commands = append(commands, constructCmd)

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
