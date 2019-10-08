package construct

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"time"
	"unicode/utf16"

	. "github.com/cloudfoundry-incubator/stembuild/remotemanager"
)

type VMConstruct struct {
	ctx             context.Context
	remoteManager   RemoteManager
	Client          IaasClient
	guestManager    GuestManager
	vmInventoryPath string
	vmUsername      string
	vmPassword      string
	winRMEnabler    WinRMEnabler
	osValidator     OSValidator
	messenger       ConstructMessenger
	poller          Poller
}

const provisionDir = "C:\\provision\\"
const stemcellAutomationName = "StemcellAutomation.zip"
const stemcellAutomationDest = provisionDir + stemcellAutomationName
const lgpoDest = provisionDir + "LGPO.zip"
const stemcellAutomationScript = provisionDir + "Setup.ps1"
const powershell = "C:\\Windows\\System32\\WindowsPowerShell\\V1.0\\powershell.exe"
const boshPsModules = "bosh-psmodules.zip"
const winRMPsScript = "BOSH.WinRM.psm1"

func NewVMConstruct(
	ctx context.Context,
	remoteManager RemoteManager,
	vmUsername,
	vmPassword,
	vmInventoryPath string,
	client IaasClient,
	guestManager GuestManager,
	winRMEnabler WinRMEnabler,
	osValidator OSValidator,
	messenger ConstructMessenger,
	poller Poller,
) *VMConstruct {

	return &VMConstruct{
		ctx,
		remoteManager,
		client,
		guestManager,
		vmInventoryPath,
		vmUsername,
		vmPassword,
		winRMEnabler,
		osValidator,
		messenger,
		poller,
	}
}

//go:generate counterfeiter . GuestManager
type GuestManager interface {
	ExitCodeForProgramInGuest(ctx context.Context, pid int64) (int32, error)
	StartProgramInGuest(ctx context.Context, command, args string) (int64, error)
	DownloadFileInGuest(ctx context.Context, path string) (io.Reader, int64, error)
}

//go:generate counterfeiter . IaasClient
type IaasClient interface {
	UploadArtifact(vmInventoryPath, artifact, destination, username, password string) error
	MakeDirectory(vmInventoryPath, path, username, password string) error
	Start(vmInventoryPath, username, password, command string, args ...string) (string, error)
	WaitForExit(vmInventoryPath, username, password, pid string) (int, error)
	IsPoweredOff(vmInventoryPath string) (bool, error)
}

//go:generate counterfeiter . WinRMEnabler
type WinRMEnabler interface {
	Enable() error
}

//go:generate counterfeiter . OSValidator
type OSValidator interface {
	Validate() error
}

//go:generate counterfeiter . ConstructMessenger
type ConstructMessenger interface {
	CreateProvisionDirStarted()
	CreateProvisionDirSucceeded()
	UploadArtifactsStarted()
	UploadArtifactsSucceeded()
	EnableWinRMStarted()
	EnableWinRMSucceeded()
	ValidateVMConnectionStarted()
	ValidateVMConnectionSucceeded()
	ExtractArtifactsStarted()
	ExtractArtifactsSucceeded()
	ExecuteScriptStarted()
	ExecuteScriptSucceeded()
	UploadFileStarted(artifact string)
	UploadFileSucceeded()
	RestartInProgress()
	ShutdownCompleted()
}

//go:generate counterfeiter . Poller
type Poller interface {
	Poll(duration time.Duration, loopFunc func() (bool, error)) error
}

func (c *VMConstruct) PrepareVM() error {

	err := c.osValidator.Validate()
	if err != nil {
		return err
	}

	err = c.createProvisionDirectory()
	if err != nil {
		return err
	}
	c.messenger.UploadArtifactsStarted()
	err = c.uploadArtifacts()
	if err != nil {
		return err
	}
	c.messenger.UploadArtifactsSucceeded()

	c.messenger.EnableWinRMStarted()
	err = c.winRMEnabler.Enable()
	if err != nil {
		return err
	}
	c.messenger.EnableWinRMSucceeded()

	c.messenger.ValidateVMConnectionStarted()
	err = c.canConnectToVM()
	if err != nil {
		return err
	}
	c.messenger.ValidateVMConnectionSucceeded()

	c.messenger.ExtractArtifactsStarted()
	err = c.extractArchive()
	if err != nil {
		return err
	}
	c.messenger.ExtractArtifactsSucceeded()

	c.messenger.ExecuteScriptStarted()
	err = c.executeSetupScript()
	if err != nil {
		return err
	}
	c.messenger.ExecuteScriptSucceeded()

	err = c.isPoweredOff(time.Minute)
	if err != nil {
		return err
	}

	return nil
}

func (c *VMConstruct) canConnectToVM() error {
	err := c.remoteManager.CanReachVM()
	if err != nil {
		return err
	}

	err = c.remoteManager.CanLoginVM()
	if err != nil {
		return err
	}

	return nil
}

func (c *VMConstruct) createProvisionDirectory() error {
	c.messenger.CreateProvisionDirStarted()
	err := c.Client.MakeDirectory(c.vmInventoryPath, provisionDir, c.vmUsername, c.vmPassword)
	if err != nil {
		return err
	}
	c.messenger.CreateProvisionDirSucceeded()
	return nil
}

func (c *VMConstruct) uploadArtifacts() error {
	c.messenger.UploadFileStarted("LGPO")
	err := c.Client.UploadArtifact(c.vmInventoryPath, "./LGPO.zip", lgpoDest, c.vmUsername, c.vmPassword)
	if err != nil {
		return err
	}
	c.messenger.UploadFileSucceeded()

	c.messenger.UploadFileStarted("stemcell preparation artifacts")
	err = c.Client.UploadArtifact(c.vmInventoryPath, fmt.Sprintf("./%s", stemcellAutomationName), stemcellAutomationDest, c.vmUsername, c.vmPassword)
	if err != nil {
		return err
	}
	c.messenger.UploadFileSucceeded()

	return nil
}

func (c *VMConstruct) extractArchive() error {
	err := c.remoteManager.ExtractArchive(stemcellAutomationDest, provisionDir)
	return err
}

func (c *VMConstruct) executeSetupScript() error {
	err := c.remoteManager.ExecuteCommand("powershell.exe " + stemcellAutomationScript)
	return err
}

func (c *VMConstruct) isPoweredOff(duration time.Duration) error {
	err := c.poller.Poll(1*time.Minute, func() (bool, error) {
		isPoweredOff, err := c.Client.IsPoweredOff(c.vmInventoryPath)

		if err != nil {
			return false, err
		}

		c.messenger.RestartInProgress()

		return isPoweredOff, nil
	})

	if err != nil {
		return err
	}

	c.messenger.ShutdownCompleted()
	return nil
}

func encodePowershellCommand(command []byte) string {
	runeCommand := []rune(string(command))
	utf16Command := utf16.Encode(runeCommand)
	byteCommand := &bytes.Buffer{}
	for _, utf16char := range utf16Command {
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, utf16char)
		byteCommand.Write(b) // This write never returns an error.
	}
	return base64.StdEncoding.EncodeToString(byteCommand.Bytes())
}
