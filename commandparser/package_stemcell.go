package commandparser

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/stembuild/filesystem"
	"github.com/cloudfoundry-incubator/stembuild/version"

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

func (*PackageCmd) Name() string { return "package" }
func (*PackageCmd) Synopsis() string {
	return "Create a BOSH Stemcell from a VMDK file or a provisioned vCenter VM"
}
func (*PackageCmd) Usage() string {
	return fmt.Sprintf(`
Create a BOSH Stemcell from a VMDK file or a provisioned vCenter VM

VM on vCenter:

  %[1]s package -vcenter-url <vCenter URL> -vcenter-username <vCenter username> -vcenter-password <vCenter password> -vm-inventory-path <vCenter VM inventory path>

  Requirements:
    - VM provisioned using the stembuild construct command
    - Access to vCenter environment
    - The [vcenter-url], [vcenter-username], [vcenter-password], and [vm-inventory-path] flags must be specified.
    - NOTE: The 'vm' keyword must be included between the datacenter name and folder name for the vm-inventory-path (e.g: <datacenter>/vm/<vm-folder>/<vm-name>) 
  Example:
    %[1]s package -vcenter-url vcenter.example.com -vcenter-username root -vcenter-password 'password' -vm-inventory-path '/my-datacenter/vm/my-folder/my-vm' 

VMDK: 

  %[1]s package -vmdk <path-to-vmdk> 

  Requirements:
    - The VMware 'ovftool' binary must be on your path or Fusion/Workstation
    must be installed (both include the 'ovftool').
    - The [vmdk] flag must be specified.  If the [output] flag is
    not specified the stemcell will be created in the current working directory.

  Example:
    %[1]s package -vmdk my-1803-vmdk.vmdk 

    Will create an Windows 1803 stemcell using [vmdk] 'my-1803-vmdk.vmdk'
    The final stemcell will be found in the current working directory.

Flags:
`, filepath.Base(os.Args[0]))
}

func (p *PackageCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.sourceConfig.Vmdk, "vmdk", "", "VMDK file to create stemcell from")
	f.StringVar(&p.sourceConfig.VmInventoryPath, "vm-inventory-path", "", "vCenter VM inventory path. (e.g: <datacenter>/vm/<vm-folder>/<vm-name>)")
	f.StringVar(&p.sourceConfig.Username, "vcenter-username", "", "vCenter username")
	f.StringVar(&p.sourceConfig.Password, "vcenter-password", "", "vCenter password")
	f.StringVar(&p.sourceConfig.URL, "vcenter-url", "", "vCenter url")
	f.StringVar(&p.outputConfig.OutputDir, "outputDir", "", "Output directory, default is the current working directory.")
	f.StringVar(&p.outputConfig.OutputDir, "o", "", "Output directory (shorthand)")
	f.StringVar(&p.sourceConfig.CaCertFile, "vcenter-ca-certs", "", "filepath for custom ca certs")
}

func (p *PackageCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	logLevel := colorlogger.NONE
	if p.GlobalFlags.Debug {
		logLevel = colorlogger.DEBUG
	}

	p.setOSandStemcellVersions()

	err := p.outputConfig.ValidateConfig()

	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		return subcommands.ExitFailure
	}

	packager, err := factory.GetPackager(p.sourceConfig, p.outputConfig, logLevel, p.GlobalFlags.Color)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return subcommands.ExitFailure
	}

	err = packager.ValidateFreeSpaceForPackage(&filesystem.OSFileSystem{})
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
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		_, _ = fmt.Fprintln(os.Stderr, "Please provide the error logs to bosh-windows-eng@pivotal.io")

		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (p *PackageCmd) setOSandStemcellVersions() {
	defaultOs, defaultStemcellVersion := version.GetVersions(version.Version)
	p.outputConfig.Os = defaultOs
	p.outputConfig.StemcellVersion = defaultStemcellVersion
}
