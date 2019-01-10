package commandparser

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/stembuild/filesystem"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/factory"

	"github.com/cloudfoundry-incubator/stembuild/colorlogger"
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/config"
	"github.com/google/subcommands"
)

type PackageCmd struct {
	GlobalFlags  *GlobalFlags
	sourceConfig config.SourceConfig
	outputConfig config.OutputConfig
}

const gigabyte = 1024 * 1024 * 1024

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
	%[1]s package -vmdk disk.vmdk -stemcell-version 1.2 -os 1803

	Will create an Windows 1803 stemcell using [vmdk] 'disk.vmdk', and set the stemcell version to 1.2.
	The final stemcell will be found in the current working directory.

Flags:
`, filepath.Base(os.Args[0]))
}

func (p *PackageCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.sourceConfig.Vmdk, "vmdk", "", "VMDK file to create stemcell from")
	f.StringVar(&p.sourceConfig.VmName, "vm-name", "", "Name of VM in vCenter")
	f.StringVar(&p.sourceConfig.Username, "username", "", "vCenter username")
	f.StringVar(&p.sourceConfig.Password, "password", "", "vCenter password")
	f.StringVar(&p.sourceConfig.URL, "url", "", "vCenter url")
	f.StringVar(&p.outputConfig.Os, "os", "", "OS version must be either 2012R2, 2016, or 1803")
	f.StringVar(&p.outputConfig.StemcellVersion, "stemcell-version", "", "Stemcell version in the form of [DIGITS].[DIGITS] (e.g. 123.01)")
	f.StringVar(&p.outputConfig.StemcellVersion, "s", "", "Stemcell version (shorthand)")
	f.StringVar(&p.outputConfig.OutputDir, "outputDir", "", "Output directory, default is the current working directory.")
	f.StringVar(&p.outputConfig.OutputDir, "o", "", "Output directory (shorthand)")
}

func (p *PackageCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	logLevel := colorlogger.NONE
	if p.GlobalFlags.Debug {
		logLevel = colorlogger.DEBUG
	}

	err := p.outputConfig.ValidateConfig()

	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		return subcommands.ExitFailure
	}

	enoughSpace, requiredSpace, err := ValidateFreeSpaceForPackage(p.sourceConfig.Vmdk, &filesystem.OSFileSystem{})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "problem checking disk space: %s", err)
		return subcommands.ExitFailure
	}
	if !enoughSpace {
		_, _ = fmt.Fprintf(os.Stderr, "Not enough space to create stemcell. Free up %d MB and try again", requiredSpace/(1024*1024))
		return subcommands.ExitFailure
	}

	packager, err := factory.GetPackager(p.sourceConfig, p.outputConfig, logLevel, p.GlobalFlags.Color)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return subcommands.ExitFailure
	}

	err = packager.ValidateSourceParameters()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return subcommands.ExitFailure
	}

	if err := packager.Package(); err != nil {
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
