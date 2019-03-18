package commandparser

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
)

//go:generate counterfeiter . VmConstruct
type VmConstruct interface {
	CanConnectToVM() error
	PrepareVM() error
}

//go:generate counterfeiter . VMPreparerFactory
type VMPreparerFactory interface {
	VMPreparer(string, string, string) VmConstruct
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
}

type ConstructCmd struct {
	winrmUsername string
	winrmPassword string
	winrmIP       string
	factory       VMPreparerFactory
	validator     ConstructCmdValidator
	messenger     ConstructMessenger
	GlobalFlags   *GlobalFlags
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
	f.StringVar(&p.winrmIP, "winrm-ip", "", "IP of machine for WinRM connection")
	f.StringVar(&p.winrmUsername, "winrm-username", "", "Username for WinRM connection")
	f.StringVar(&p.winrmPassword, "winrm-password", "", "Password for WinRM connection. Needs to be wrapped in single quotations.")
}

func (p *ConstructCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if !p.validator.PopulatedArgs(p.winrmIP, p.winrmUsername, p.winrmPassword) {
		p.messenger.ArgumentsNotProvided()
		return subcommands.ExitFailure
	}
	if !p.validator.LGPOInDirectory() {
		p.messenger.LGPONotFound()
		return subcommands.ExitFailure
	}

	vmConstruct := p.factory.VMPreparer(p.winrmIP, p.winrmUsername, p.winrmPassword)
	err := vmConstruct.PrepareVM()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
