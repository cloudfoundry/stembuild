package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
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
	return `package -vmdk <path-to-vmdk> -os <os version> -version <NOT SURE WHAT THIS IS THE VERSION OF???>`
}

func (p *packageCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.vmdk, "vmdk", "", "vmdk to convert to stemcell")
	f.StringVar(&p.os, "os", "", "os to build for, e.g. 2012R2, 1709, 1803")
	f.StringVar(&p.version, "version", "", "again, DONT KNOW WHAT TO PUT HERE")
}

func (p *packageCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

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
