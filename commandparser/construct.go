package commandparser

import (
	"context"
	"flag"
	"fmt"
	"github.com/cloudfoundry-incubator/stembuild/colorlogger"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
)

//go:generate counterfeiter . VMPreparer
type VMPreparer interface {
	PrepareVM() error
}

//go:generate counterfeiter . VMPreparerFactory
type VMPreparerFactory interface {
	VMPreparer(string, string, string) VMPreparer
}

//go:generate counterfeiter . ConstructCmdValidator
type ConstructCmdValidator interface {
	NonEmptyArgs(...string) bool
	LGPOInDirectory() bool
	ValidStemcellInfo(string) bool
}

//go:generate counterfeiter . ConstructMessenger
type ConstructMessenger interface {
	ArgumentsNotProvided() error
	InvalidStemcellVersion(string) error
	LGPONotFound() error
}

type ConstructCmd struct {
	stemcellVersion string
	winrmUsername   string
	winrmPassword   string
	winrmIP         string
	factory         VMPreparerFactory
	validator       ConstructCmdValidator
	GlobalFlags     *GlobalFlags
}

func NewConstructCmd(factory VMPreparerFactory, validator ConstructCmdValidator) ConstructCmd {
	return ConstructCmd{factory: factory, validator: validator}
}

func (*ConstructCmd) Name() string { return "construct" }
func (*ConstructCmd) Synopsis() string {
	return "Provisions and syspreps an existing VM on vCenter, ready to be packaged into a stemcell"
}

func (*ConstructCmd) Usage() string {
	return fmt.Sprintf(`%[1]s construct -stemcell-version <stemcell version> -winrm-ip <IP of VM> -winrum-username <WinRm username> -winrm-password <WinRm password>

Prepares a VM to be used by stembuild package. It leverages stemcell automation scripts to provision a VM to be used as a stemcell.

Requirements:
	LGPO.zip in current working directory
	Running Windows VM with:
		- Up to date Operating System
		- WinRm enabled
		- Reachable by IP
		- Username and password with Administrator privileges
	The [stemcell-version], [ip], [winrm-username], [winrm-password] flags must be specified.

Example:
	%[1]s construct -stemcell-version 1803.1 -winrm-ip '10.0.0.5' -winrm-username Admin -winrm-password 'password'

Flags:
`, filepath.Base(os.Args[0]))
}

func (p *ConstructCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.stemcellVersion, "stemcell-version", "", "Stemcell version in the form of [DIGITS].[DIGITS] (e.g. 1803.1)")
	f.StringVar(&p.winrmIP, "winrm-ip", "", "IP of machine for WinRM connection")
	f.StringVar(&p.winrmUsername, "winrm-username", "", "Username for WinRM connection")
	f.StringVar(&p.winrmPassword, "winrm-password", "", "Password for WinRM connection. Needs to be wrapped in single quotations.")
}

func (p *ConstructCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	logLevel := colorlogger.NONE
	if p.GlobalFlags.Debug {
		logLevel = colorlogger.DEBUG
	}
	logger := colorlogger.ConstructLogger(logLevel, p.GlobalFlags.Color, os.Stderr)
	logger.Debugf("hello, world.")

	if !p.validator.NonEmptyArgs(p.winrmIP, p.winrmUsername, p.winrmPassword, p.stemcellVersion) {
		_, _ = fmt.Fprintf(os.Stderr, "All arguments must be provided")
		return subcommands.ExitFailure
	}
	if !p.validator.ValidStemcellInfo(p.stemcellVersion) {
		_, _ = fmt.Fprintf(os.Stderr, "invalid stemcellVersion (%s) expected format [NUMBER].[NUMBER] or "+
			"[NUMBER].[NUMBER].[NUMBER]\n", p.stemcellVersion)
		return subcommands.ExitFailure
	}
	if !p.validator.LGPOInDirectory() {
		_, _ = fmt.Fprintf(os.Stderr, "lgpo not found in current directory")
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
