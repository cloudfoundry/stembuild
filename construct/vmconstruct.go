package construct

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/cloudfoundry-incubator/stembuild/assets"
	. "github.com/cloudfoundry-incubator/stembuild/remotemanager"
	"unicode/utf16"
)

type VMConstruct struct {
	remoteManager   RemoteManager
	Client          IaasClient
	vmInventoryPath string
	vmUsername      string
	vmPassword      string
	unarchiver      zipUnarchiver
}

const provisionDir = "C:\\provision\\"
const stemcellAutomationDest = provisionDir + "StemcellAutomation.zip"
const lgpoDest = provisionDir + "LGPO.zip"
const stemcellAutomationScript = provisionDir + "Setup.ps1"
const powershell = "C:\\Windows\\System32\\WindowsPowerShell\\V1.0\\powershell.exe"

func NewVMConstruct(winrmIP, winrmUsername, winrmPassword, vmInventoryPath string, client IaasClient, unarchiver zipUnarchiver) *VMConstruct {
	return &VMConstruct{NewWinRM(winrmIP, winrmUsername, winrmPassword), client, vmInventoryPath, winrmUsername, winrmPassword, unarchiver}
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

func (c *VMConstruct) CanConnectToVM() error {
	fmt.Print("Validating connection to vm...")
	err := c.remoteManager.CanReachVM()
	if err != nil {
		return err
	}
	err = c.remoteManager.CanLoginVM()
	if err != nil {
		return err
	}
	fmt.Println(" succeeded.")

	return nil
}

func (c *VMConstruct) PrepareVM() error {
	fmt.Println("\nTransferring ~20 MB to the Windows VM. Depending on your connection, the transfer may take 15-45 minutes")
	err := c.uploadArtifacts()
	if err != nil {
		return err
	}
	fmt.Println("All files have been uploaded.")

	fmt.Print("\nExtracting artifacts...")
	err = c.extractArchive()
	if err != nil {
		return err
	}
	fmt.Println(" Artifacts Extracted.")

	fmt.Println("\nExecuting setup script...")
	err = c.executeSetupScript()
	if err != nil {
		return err
	}
	fmt.Println("\nFinished executing setup script.")

	return nil
}

func (c *VMConstruct) uploadArtifacts() error {

	fmt.Print("\tCreating provision dir on target VM...")
	err := c.Client.MakeDirectory(c.vmInventoryPath, provisionDir, c.vmUsername, c.vmPassword)
	if err != nil {
		return err
	}
	fmt.Println(" Finished creating provision dir.")

	fmt.Print("\tUploading LGPO to target VM...")
	err = c.Client.UploadArtifact(c.vmInventoryPath, "./LGPO.zip", lgpoDest, c.vmUsername, c.vmPassword)
	if err != nil {
		return err
	}
	fmt.Println(" Finished uploading LGPO.")

	fmt.Print("\tUploading stemcell preparation artifacts to target VM...")
	err = c.Client.UploadArtifact(c.vmInventoryPath, "./StemcellAutomation.zip", stemcellAutomationDest, c.vmUsername, c.vmPassword)
	if err != nil {
		return err
	}
	fmt.Println(" Finished uploading artifacts to target VM.")

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

	saZip, err := assets.Asset("StemcellAutomation.zip")
	if err != nil {
		return fmt.Errorf(failureString, err)
	}

	bmZip, err := c.unarchiver.Unzip(saZip, "bosh-modules.zip")
	if err != nil {
		return fmt.Errorf(failureString, err)
	}

	rawWinRM, err := c.unarchiver.Unzip(bmZip, "BOSH.WinRM.psm1")
	if err != nil {
		return fmt.Errorf(failureString, err)
	}

	// Since BOSH.WinRM.psm1 just contains the enable WinRM function, we need to append 'Enable-WinRM' in order
	// for the function to be executed.
	rawWinRMwtCmd := append(rawWinRM, []byte("\nEnable-WinRM\n")...)

	//TODO: Maybe extract this code block
	runeWinRM := []rune(string(rawWinRMwtCmd))
	utf16WinRM := utf16.Encode(runeWinRM)
	byteWinRM := &bytes.Buffer{}
	for _, utf16char := range utf16WinRM {
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, utf16char)
		byteWinRM.Write(b)
	}
	base64WinRM := base64.StdEncoding.EncodeToString(byteWinRM.Bytes())

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

	fmt.Println("WinRm enabled on the guest VM")

	return nil
}
