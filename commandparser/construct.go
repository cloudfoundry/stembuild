package commandparser

import (
	"context"
	"flag"
	"fmt"
	. "github.com/cloudfoundry-incubator/stembuild/construct"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/stembuild/colorlogger"
	"github.com/google/subcommands"
)

type ConstructCmd struct {
	stemcellVersion string
	winrmUsername   string
	winrmPassword   string
	winrmIP         string
	GlobalFlags     *GlobalFlags
}

func (*ConstructCmd) Name() string     { return "construct" }
func (*ConstructCmd) Synopsis() string { return "Transfer automation artifact and LGPO to a vCenter VM" }

func (*ConstructCmd) Usage() string {
	return fmt.Sprintf(`%[1]s construct -stemcell-version <stemcell stemcellVersion> -winrm-ip <winrm ip of VM> -winrm-username <winrm username> -winrm-password <winrm password>

Prepare a VM to be used by stembuild package. It leverages Stemcell Automation Scripts to construct a VM to be used as a stemcell.

The [stemcell-version], [winrm-ip], [winrm-username], [winrm-password] flags must be specified.

Requirements:
	Running Windows VM with:
		- Up to date Operating System
		- WinRm enabled
		- Reachable by IP
		- Username and password with Administrator privileges 
	StemcellAutomation.zip in current working directory
	LGPO.zip in current working directory

Examples:
	%[1]s construct -stemcell-version 1709.1 -winrm-ip '10.0.0.5' -winrm-username Admin -winrm-password 'password'

	This will connect to VM with IP 10.0.0.5 using credentials Admin:password, upload and execute StemcellAutomation.zip found in the working directory.
	StemcellAutomation.zip requires LGPO.zip to be present in the working directory.
	When command exits successfully, the VM will be Sysprepped and powered off.

Flags:
`, filepath.Base(os.Args[0]))
}

func (p *ConstructCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.stemcellVersion, "stemcell-version", "", "Stemcell version in the form of [DIGITS].[DIGITS] (e.g. 123.01)")
	f.StringVar(&p.stemcellVersion, "s", "", "Stemcell version (shorthand)")
	f.StringVar(&p.winrmIP, "winrm-ip", "", "IP of machine for WinRM connection")
	f.StringVar(&p.winrmIP, "ip", "", "winrm-ip (shorthand)")
	f.StringVar(&p.winrmUsername, "winrm-username", "", "Username for winRM connection")
	f.StringVar(&p.winrmUsername, "u", "", "winrm-username (shorthand)")
	f.StringVar(&p.winrmPassword, "winrm-password", "", "Password for winRM connection")
	f.StringVar(&p.winrmPassword, "p", "", "winrm-password (shorthand)")
}

func (p *ConstructCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	logLevel := colorlogger.NONE
	if p.GlobalFlags.Debug {
		logLevel = colorlogger.DEBUG
	}
	logger := colorlogger.ConstructLogger(logLevel, p.GlobalFlags.Color, os.Stderr)
	logger.Debugf("hello, world.")
	if !IsValidStemcellVersion(p.stemcellVersion) {
		_, _ = fmt.Fprintf(os.Stderr, "invalid stemcellVersion (%s) expected format [NUMBER].[NUMBER] or "+
			"[NUMBER].[NUMBER].[NUMBER]\n", p.stemcellVersion)

		return subcommands.ExitFailure
	}

	pwd, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "unable to find current working directory", err)
		return subcommands.ExitFailure
	}
	automationArtifactPresent, err := IsArtifactInDirectory(pwd, "StemcellAutomation.zip")
	if !automationArtifactPresent {
		_, _ = fmt.Fprintf(os.Stderr, "automation artifact not found in current directory")
		return subcommands.ExitFailure
	}
	lgpoPresent, err := IsArtifactInDirectory(pwd, "LGPO.zip")
	if !lgpoPresent {
		_, _ = fmt.Fprintf(os.Stderr, "lgpo not found in current directory")
		return subcommands.ExitFailure
	}

	vmConstruct := NewVMConstruct(p.winrmIP, p.winrmUsername, p.winrmPassword)
	err = vmConstruct.PrepareVM()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
