package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/pivotal-cf-experimental/stembuild/ovftool"
	"github.com/pivotal-cf-experimental/stembuild/pack"
	. "github.com/pivotal-cf-experimental/stembuild/pack/options"
	"github.com/pivotal-cf-experimental/stembuild/stemcell"
	"path/filepath"

	"log"
	"os"
	"path"
	"strings"
)

type packageCmd struct {
	vmdk      string
	os        string
	version   string
	outputDir string
}

type globalFlags struct {
	debug bool
	color bool
}

var gf globalFlags

func (*packageCmd) Name() string     { return "package" }
func (*packageCmd) Synopsis() string { return "ADD A GOOD SYNOPSIS" }
func (*packageCmd) Usage() string {
	return "package --vmdk <path-to-vmdk> -os <os version> -version <stemcell version>\n"
}

func (p *packageCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.vmdk, "vmdk", "", "VMDK file to create stemcell from")
	f.StringVar(&p.os, "os", "", "OS version must be either 2012R2, 2016 or 1803")
	f.StringVar(&p.version, "version", "", "Stemcell version in the form of [DIGITS].[DIGITS] (e.g. 123.01)")
	f.StringVar(&p.version, "v", "", "Stemcell version (shorthand)")
	f.StringVar(&p.outputDir, "outputDir", "", "Output directory, default is the current working directory.")
	f.StringVar(&p.outputDir, "o", "", "Output directory (shorthand)")
}
func (g *globalFlags) getDebug() func(format string, a ...interface{}) {

	debugFunc := func(format string, a ...interface{}) {}
	prefix := "debug: "
	if g.color {
		prefix = "\033[32m" + prefix + "\033[0m"
	}
	if g.debug {
		debugFunc = log.New(os.Stderr, prefix, 0).Printf
	}
	return debugFunc
}
func (p *packageCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	if validVMDK, err := pack.IsValidVMDK(p.vmdk); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return subcommands.ExitFailure
	} else if !validVMDK {
		fmt.Fprintf(os.Stderr, "VMDK not specified or invalid\n")
		return subcommands.ExitFailure
	}
	if !pack.IsValidOS(p.os) {
		fmt.Fprintf(os.Stderr, "OS version must be either 2012R2, 1709, or 1803 have: %s\n", p.os)
		return subcommands.ExitFailure
	}
	if !pack.IsValidVersion(p.version) {
		fmt.Fprintf(os.Stderr, "invalid version (%s) expected format [NUMBER].[NUMBER] or "+
			"[NUMBER].[NUMBER].[NUMBER]\n", p.version)

		return subcommands.ExitFailure
	}

	if p.outputDir == "" || p.outputDir == "." {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting working directory %s", err)
			return subcommands.ExitFailure
		}
		p.outputDir = cwd
	} else if err := pack.ValidateOrCreateOutputDir(p.outputDir); err != nil {
		return subcommands.ExitFailure
	}

	name := filepath.Join(p.outputDir, stemcell.StemcellFilename(p.version, p.os))
	gf.getDebug()("validating that stemcell filename (%s) does not exist", name)
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error with output file (%s): %v (file may already exist)", name, err)
		return subcommands.ExitFailure
	}

	fmt.Print("Finding 'ovftool'...")
	searchPaths, err := ovftool.SearchPaths()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get search paths for Ovftool: %s", err)
		return subcommands.ExitFailure
	}
	path, err := ovftool.Ovftool(searchPaths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not locate 'ovftool' on PATH: %s", err)
		return subcommands.ExitFailure
	}
	fmt.Printf("...'ovftool' found at: %s\n", path)

	c := stemcell.Config{
		Stop:         make(chan struct{}),
		Debugf:       gf.getDebug(),
		BuildOptions: StembuildOptions{},
	}

	c.BuildOptions.VMDKFile = p.vmdk
	c.BuildOptions.OSVersion = strings.ToUpper(p.os)
	c.BuildOptions.Version = p.version
	c.BuildOptions.OutputDir = p.outputDir

	if err := c.Package(); err != nil {
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func main() {

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.BoolVar(&gf.debug, "debug", false, "Print lots of debugging informatio")
	fs.BoolVar(&gf.color, "color", false, "Colorize debug output")

	commander := subcommands.NewCommander(fs, path.Base(os.Args[0]))

	commander.Register(commander.HelpCommand(), "")
	commander.Register(commander.FlagsCommand(), "")
	commander.Register(commander.CommandsCommand(), "")

	commander.Register(&packageCmd{}, "")

	fs.Parse(os.Args[1:])
	ctx := context.Background()
	os.Exit(int(commander.Execute(ctx)))
}
