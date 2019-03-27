package construct

import (
	"fmt"
	. "github.com/cloudfoundry-incubator/stembuild/remotemanager"
)

type VMConstruct struct {
	remoteManager   RemoteManager
	Client          IaasClient
	vmInventoryPath string
	vmUsername      string
	vmPassword      string
}

const provisionDir = "C:\\provision\\"
const stemcellAutomationDest = provisionDir + "StemcellAutomation.zip"
const lgpoDest = provisionDir + "LGPO.zip"
const stemcellAutomationScript = provisionDir + "Setup.ps1"

func NewVMConstruct(winrmIP, winrmUsername, winrmPassword, vmInventoryPath string, client IaasClient) *VMConstruct {
	return &VMConstruct{NewWinRM(winrmIP, winrmUsername, winrmPassword), client, vmInventoryPath, winrmUsername, winrmPassword}
}

//go:generate counterfeiter . IaasClient
type IaasClient interface {
	UploadArtifact(vmInventoryPath, artifact, destination, username, password string) error
	MakeDirectory(vmInventoryPath, path, username, password string) error
	Start(vmInventoryPath, username, password, command string, args ...string) (string, error)
	WaitForExit(vmInventoryPath, username, password, pid string) (int, error)
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
