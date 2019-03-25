package commandparser

import (
	"context"
	"flag"
	"fmt"
	"github.com/cloudfoundry-incubator/stembuild/construct/config"
	"github.com/google/subcommands"
	"os"
	"path/filepath"
)

//go:generate counterfeiter . VmConstruct
type VmConstruct interface {
	CanConnectToVM() error
	PrepareVM() error
}

//go:generate counterfeiter . VMPreparerFactory
type VMPreparerFactory interface {
	VMPreparer(config config.SourceConfig) VmConstruct
}

//go:generate counterfeiter . ConstructCmdValidator
type ConstructCmdValidator interface {
	PopulatedArgs(...string) bool
	LGPOInDirectory() bool
}

//go:generate counterfeiter . ConstructMessenger
type ConstructMessenger interface {
	ArgumentsNotProvided()
	LGPONotFound()
	CannotConnectToVM(err error)
	CannotPrepareVM(err error)
}

type ConstructCmd struct {
	sourceConfig config.SourceConfig
	factory      VMPreparerFactory
	validator    ConstructCmdValidator
	messenger    ConstructMessenger
	GlobalFlags  *GlobalFlags
}

func NewConstructCmd(factory VMPreparerFactory, validator ConstructCmdValidator, messenger ConstructMessenger) ConstructCmd {
	return ConstructCmd{factory: factory, validator: validator, messenger: messenger}
}

func (*ConstructCmd) Name() string { return "construct" }
func (*ConstructCmd) Synopsis() string {
	return "Provisions and syspreps an existing VM on vCenter, ready to be packaged into a stemcell"
}

func (*ConstructCmd) Usage() string {
	return fmt.Sprintf(`%[1]s construct -winrm-ip <IP of VM> -winrum-username <WinRm username> -winrm-password <WinRm password>

Prepares a VM to be used by stembuild package. It leverages stemcell automation scripts to provision a VM to be used as a stemcell.

Requirements:
	LGPO.zip in current working directory
	Running Windows VM with:
		- Up to date Operating System
		- WinRm enabled
		- Reachable by IP
		- Username and password with Administrator privileges
	The [winrm-ip], [winrm-username], [winrm-password] flags must be specified.

Example:
	%[1]s construct -winrm-ip '10.0.0.5' -winrm-username Admin -winrm-password 'password'

Flags:
`, filepath.Base(os.Args[0]))
}

func (p *ConstructCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.sourceConfig.GuestVmIp, "winrm-ip", "", "IP of machine for WinRM connection")
	f.StringVar(&p.sourceConfig.GuestVMUsername, "winrm-username", "", "Username for WinRM connection")
	f.StringVar(&p.sourceConfig.GuestVMPassword, "winrm-password", "", "Password for WinRM connection. Needs to be wrapped in single quotations.")
	f.StringVar(&p.sourceConfig.VCenterUrl, "vcenter-url", "", "vCenter url")
	f.StringVar(&p.sourceConfig.VCenterUsername, "vcenter-username", "", "vCenter username")
	f.StringVar(&p.sourceConfig.VCenterPassword, "vcenter-password", "", "vCenter password")
	f.StringVar(&p.sourceConfig.VmInventoryPath, "vm-inventory-path", "", "vCenter VM inventory path. (e.g: <datacenter>/vm/<vm-folder>/<vm-name>)")
}

func (p *ConstructCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	c := p.sourceConfig
	if !p.validator.PopulatedArgs(c.GuestVmIp, c.GuestVMUsername, c.GuestVMPassword, c.VCenterUrl, c.VCenterUsername, c.VCenterPassword, c.VmInventoryPath) {
		p.messenger.ArgumentsNotProvided()
		return subcommands.ExitFailure
	}
	if !p.validator.LGPOInDirectory() {
		p.messenger.LGPONotFound()
		return subcommands.ExitFailure
	}

	vmConstruct := p.factory.VMPreparer(p.sourceConfig)

	err := vmConstruct.CanConnectToVM()
	if err != nil {
		p.messenger.CannotConnectToVM(err)
		return subcommands.ExitFailure
	}

	err = vmConstruct.PrepareVM()
	if err != nil {
		p.messenger.CannotPrepareVM(err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
