package commandparser

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	. "github.com/pivotal-cf-experimental/stembuild/pack/options"
	"github.com/pivotal-cf-experimental/stembuild/pack/ovftool"
	"github.com/pivotal-cf-experimental/stembuild/pack/stemcell"
	"os"
	"path/filepath"
	"strings"
)

type PackageCmd struct {
	vmdk        string
	os          string
	version     string
	outputDir   string
	GlobalFlags *GlobalFlags
}

func (*PackageCmd) Name() string     { return "package" }
func (*PackageCmd) Synopsis() string { return "ADD A GOOD SYNOPSIS" }
func (*PackageCmd) Usage() string {
	return "package --vmdk <path-to-vmdk> -os <os version> -version <stemcell version>\n"
}

func (p *PackageCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.vmdk, "vmdk", "", "VMDK file to create stemcell from")
	f.StringVar(&p.os, "os", "", "OS version must be either 2012R2, 2016 or 1803")
	f.StringVar(&p.version, "version", "", "Stemcell version in the form of [DIGITS].[DIGITS] (e.g. 123.01)")
	f.StringVar(&p.version, "v", "", "Stemcell version (shorthand)")
	f.StringVar(&p.outputDir, "outputDir", "", "Output directory, default is the current working directory.")
	f.StringVar(&p.outputDir, "o", "", "Output directory (shorthand)")
}
func (p *PackageCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	if validVMDK, err := IsValidVMDK(p.vmdk); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return subcommands.ExitFailure
	} else if !validVMDK {
		fmt.Fprintf(os.Stderr, "VMDK not specified or invalid\n")
		return subcommands.ExitFailure
	}
	if !IsValidOS(p.os) {
		fmt.Fprintf(os.Stderr, "OS version must be either 2012R2, 1709, or 1803 have: %s\n", p.os)
		return subcommands.ExitFailure
	}
	if !IsValidVersion(p.version) {
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
	} else if err := ValidateOrCreateOutputDir(p.outputDir); err != nil {
		return subcommands.ExitFailure
	}

	name := filepath.Join(p.outputDir, stemcell.StemcellFilename(p.version, p.os))
	p.GlobalFlags.GetDebug()("validating that stemcell filename (%s) does not exist", name)
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
	ovfPath, err := ovftool.Ovftool(searchPaths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not locate 'ovftool' on PATH: %s", err)
		return subcommands.ExitFailure
	}
	fmt.Printf("...'ovftool' found at: %s\n", ovfPath)

	c := stemcell.Config{
		Stop:         make(chan struct{}),
		Debugf:       p.GlobalFlags.GetDebug(),
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
