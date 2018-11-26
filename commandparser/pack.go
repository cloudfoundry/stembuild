package commandparser

import (
	"context"
	"flag"
	"fmt"
	"github.com/cloudfoundry-incubator/stembuild/colorlogger"
	. "github.com/cloudfoundry-incubator/stembuild/pack/options"
	"github.com/cloudfoundry-incubator/stembuild/pack/ovftool"
	"github.com/cloudfoundry-incubator/stembuild/pack/stemcell"
	"github.com/google/subcommands"
	"os"
	"path/filepath"
	"strings"
)

type PackageCmd struct {
	vmdk            string
	os              string
	stemcellVersion string
	outputDir       string
	GlobalFlags     *GlobalFlags
}

func (*PackageCmd) Name() string     { return "package" }
func (*PackageCmd) Synopsis() string { return "Create a BOSH Stemcell from a VMDK file" }
func (*PackageCmd) Usage() string {
	return fmt.Sprintf(`%[1]s package -vmdk <path-to-vmdk> -stemcellVersion <stemcell stemcellVersion> -os <os stemcellVersion>

Create a BOSH Stemcell from a VMDK file

The [vmdk], [stemcellVersion], and [os] flags must be specified.  If the [output] flag is
not specified the stemcell will be created in the current working directory.

Requirements:
	The VMware 'ovftool' binary must be on your path or Fusion/Workstation
	must be installed (both include the 'ovftool').

Examples:
	%[1]s -vmdk disk.vmdk -stemcellVersion 1.2 -os 1803

	Will create an Windows 1803 stemcell using [vmdk] 'disk.vmdk', and set the stemcell version to 1.2.
	The final stemcell will be found in the current working directory.

Flags:
`, filepath.Base(os.Args[0]))
}

func (p *PackageCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.vmdk, "vmdk", "", "VMDK file to create stemcell from")
	f.StringVar(&p.os, "os", "", "OS version must be either 2012R2, 2016, or 1803")
	f.StringVar(&p.stemcellVersion, "stemcell-version", "", "Stemcell version in the form of [DIGITS].[DIGITS] (e.g. 123.01)")
	f.StringVar(&p.stemcellVersion, "s", "", "Stemcell version (shorthand)")
	f.StringVar(&p.outputDir, "outputDir", "", "Output directory, default is the current working directory.")
	f.StringVar(&p.outputDir, "o", "", "Output directory (shorthand)")
}
func (p *PackageCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	logLevel := colorlogger.NONE
	if p.GlobalFlags.Debug {
		logLevel = colorlogger.DEBUG
	}
	logger := colorlogger.ConstructLogger(logLevel, p.GlobalFlags.Color, os.Stderr)

	if validVMDK, err := IsValidVMDK(p.vmdk); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return subcommands.ExitFailure
	} else if !validVMDK {
		_, _ = fmt.Fprintf(os.Stderr, "VMDK not specified or invalid\n")
		return subcommands.ExitFailure
	}
	if !IsValidOS(p.os) {
		_, _ = fmt.Fprintf(os.Stderr, "OS version must be either 2012R2, 2016, or 1803 have: %s\n", p.os)
		return subcommands.ExitFailure
	}
	if !IsValidStemcellVersion(p.stemcellVersion) {
		_, _ = fmt.Fprintf(os.Stderr, "invalid stemcell version (%s) expected format [NUMBER].[NUMBER] or "+
			"[NUMBER].[NUMBER].[NUMBER]\n", p.stemcellVersion)

		return subcommands.ExitFailure
	}

	if p.outputDir == "" || p.outputDir == "." {
		cwd, err := os.Getwd()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error getting working directory %s", err)
			return subcommands.ExitFailure
		}
		p.outputDir = cwd
	} else if err := ValidateOrCreateOutputDir(p.outputDir); err != nil {
		return subcommands.ExitFailure
	}

	name := filepath.Join(p.outputDir, stemcell.StemcellFilename(p.stemcellVersion, p.os))
	logger.Debugf("validating that stemcell filename (%s) does not exist", name)
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		_, _ = fmt.Fprintf(os.Stderr, "error with output file (%s): %v (file may already exist)", name, err)
		return subcommands.ExitFailure
	}

	fmt.Print("Finding 'ovftool'...")
	searchPaths, err := ovftool.SearchPaths()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not get search paths for Ovftool: %s", err)
		return subcommands.ExitFailure
	}
	ovfPath, err := ovftool.Ovftool(searchPaths)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not locate 'ovftool' on PATH: %s", err)
		return subcommands.ExitFailure
	}
	fmt.Printf("...'ovftool' found at: %s\n", ovfPath)

	c := stemcell.Config{
		Stop:         make(chan struct{}),
		Debugf:       logger.Debugf,
		BuildOptions: StembuildOptions{},
	}

	c.BuildOptions.VMDKFile = p.vmdk
	c.BuildOptions.OSVersion = strings.ToUpper(p.os)
	c.BuildOptions.Version = p.stemcellVersion
	c.BuildOptions.OutputDir = p.outputDir

	if err := c.Package(); err != nil {
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
