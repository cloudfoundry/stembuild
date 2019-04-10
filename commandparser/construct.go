package commandparser

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/stembuild/construct/config"
	"github.com/google/subcommands"
)

//go:generate counterfeiter . VmConstruct
type VmConstruct interface {
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
	return fmt.Sprintf(`%[1]s construct -vm-ip <IP of VM> -vm-username <vm username> -vm-password <vm password>  -vcenter-url <vCenter URL> -vcenter-username <vCenter username> -vcenter-password <vCenter password> -vm-inventory-path <vCenter VM inventory path>

Prepares a VM to be used by stembuild package. It leverages stemcell automation scripts to provision a VM to be used as a stemcell.

Requirements:
	LGPO.zip in current working directory
	Running Windows VM with:
		- Up to date Operating System
		- Reachable by IP
		- Username and password with Administrator privileges
		- vCenter URL, username and password
		- vCenter Inventory Path
	The [vm-ip], [vm-username], [vm-password], [vcenter-url], [vcenter-username], [vcenter-password], [vm-inventory-path] must be specified

Example:
	%[1]s construct -vm-ip '10.0.0.5' -vm-username Admin -vm-password 'password' -vcenter-url vcenter.example.com -vcenter-username root -vcenter-password 'password' -vm-inventory-path '/datacenter/vm/folder/vm-name'

Flags:
`, filepath.Base(os.Args[0]))
}

func (p *ConstructCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.sourceConfig.GuestVmIp, "vm-ip", "", "IP of target machine")
	f.StringVar(&p.sourceConfig.GuestVMUsername, "vm-username", "", "Username of target machine")
	f.StringVar(&p.sourceConfig.GuestVMPassword, "vm-password", "", "Password of target machine. Needs to be wrapped in single quotations.")
	f.StringVar(&p.sourceConfig.VCenterUrl, "vcenter-url", "", "vCenter url")
	f.StringVar(&p.sourceConfig.VCenterUsername, "vcenter-username", "", "vCenter username")
	f.StringVar(&p.sourceConfig.VCenterPassword, "vcenter-password", "", "vCenter password")
	f.StringVar(&p.sourceConfig.VmInventoryPath, "vm-inventory-path", "", "vCenter VM inventory path. (e.g: <datacenter>/vm/<vm-folder>/<vm-name>)")
	f.StringVar(&p.sourceConfig.CaCertFile, "vcenter-ca-certs", "", "filepath for custom ca certs")
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

	err := vmConstruct.PrepareVM()
	if err != nil {
		p.messenger.CannotPrepareVM(err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
