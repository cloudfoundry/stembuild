package construct

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"unicode/utf16"

	"github.com/cloudfoundry-incubator/stembuild/assets"
	. "github.com/cloudfoundry-incubator/stembuild/remotemanager"
)

type VMConstruct struct {
	remoteManager   RemoteManager
	Client          IaasClient
	vmInventoryPath string
	vmUsername      string
	vmPassword      string
	unarchiver      zipUnarchiver
	messenger       ConstructMessenger
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
	vmIP,
	vmUsername,
	vmPassword,
	vmInventoryPath string,
	client IaasClient,
	unarchiver zipUnarchiver,
	messenger ConstructMessenger,
) *VMConstruct {

	return &VMConstruct{
		NewWinRM(vmIP, vmUsername, vmPassword),
		client,
		vmInventoryPath,
		vmUsername,
		vmPassword,
		unarchiver,
		messenger,
	}
}

//go:generate counterfeiter . IaasClient
type IaasClient interface {
	UploadArtifact(vmInventoryPath, artifact, destination, username, password string) error
	MakeDirectory(vmInventoryPath, path, username, password string) error
	Start(vmInventoryPath, username, password, command string, args ...string) (string, error)
	WaitForExit(vmInventoryPath, username, password, pid string) (int, error)
}

//go:generate counterfeiter . zipUnarchiver
type zipUnarchiver interface {
	Unzip(fileArchive []byte, file string) ([]byte, error)
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
}

func (c *VMConstruct) PrepareVM() error {
	err := c.createProvisionDirectory()
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
	err = c.enableWinRM()
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

func (c *VMConstruct) enableWinRM() error {
	failureString := "failed to enable WinRM: %s"

	saZip, err := assets.Asset(stemcellAutomationName)
	if err != nil {
		return fmt.Errorf(failureString, err)
	}

	bmZip, err := c.unarchiver.Unzip(saZip, boshPsModules)
	if err != nil {
		return fmt.Errorf(failureString, err)
	}

	rawWinRM, err := c.unarchiver.Unzip(bmZip, winRMPsScript)
	if err != nil {
		return fmt.Errorf(failureString, err)
	}

	// Since BOSH.WinRM.psm1 just contains the enable WinRM function, we need to append 'Enable-WinRM' in order
	// for the function to be executed.
	rawWinRMwtCmd := append(rawWinRM, []byte("\nEnable-WinRM\n")...)

	base64WinRM := encodePowershellCommand(rawWinRMwtCmd)

	pid, err := c.Client.Start(c.vmInventoryPath, c.vmUsername, c.vmPassword, powershell, "-EncodedCommand", base64WinRM)
	if err != nil {
		return fmt.Errorf(failureString, err)
	}

	exitCode, err := c.Client.WaitForExit(c.vmInventoryPath, c.vmUsername, c.vmPassword, pid)
	if err != nil {
		return fmt.Errorf(failureString, err)
	}
	if exitCode != 0 {
		return fmt.Errorf(failureString, fmt.Sprintf("WinRM process on guest VM exited with code %d", exitCode))
	}

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
