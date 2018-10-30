package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/pivotal-cf-experimental/stembuild"
	"os"
)

type packageCmd struct {
	vmdk    string
	os      string
	version string
}

func (*packageCmd) Name() string     { return "package" }
func (*packageCmd) Synopsis() string { return "ADD A GOOD SYNOPSIS" }
func (*packageCmd) Usage() string {
	return `package -vmdk <path-to-vmdk> -os <os version> -version <stemcell version>`
}

func (p *packageCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.vmdk, "vmdk", "", "vmdk to convert to stemcell")
	f.StringVar(&p.os, "os", "", "os to build for, e.g. 2012R2, 1709, 1803")
	f.StringVar(&p.version, "version", "", "stemcell version, e.g. 1803.7")
}

func (p *packageCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	if validVMDK, err := stembuild.IsValidVMDK(p.vmdk); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return subcommands.ExitFailure
	} else if !validVMDK {
		fmt.Fprintf(os.Stderr, "VMDK not specified or invalid\n")
		return subcommands.ExitFailure
	}
	if !stembuild.IsValidOS(p.os) {
		fmt.Fprintf(os.Stderr, "OS version must be either 2012R2, 1709, or 1803 have: %s\n", p.os)
		return subcommands.ExitFailure
	}
	if !stembuild.IsValidVersion(p.version) {
		fmt.Fprintf(os.Stderr, "invalid version (%s) expected format [NUMBER].[NUMBER] or "+
			"[NUMBER].[NUMBER].[NUMBER]\n", p.version)

		return subcommands.ExitFailure
	}

	fmt.Println("just makin' sure we're here")
	return subcommands.ExitSuccess
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&packageCmd{}, "Custom")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
