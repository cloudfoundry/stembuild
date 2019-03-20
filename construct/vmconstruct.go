package construct

import (
	"fmt"
	. "github.com/cloudfoundry-incubator/stembuild/remotemanager"
)

type VMConstruct struct {
	remoteManager RemoteManager
}

const provisionDir = "C:\\provision\\"
const stemcellAutomationDest = provisionDir + "StemcellAutomation.zip"
const lgpoDest = provisionDir + "LGPO.zip"
const stemcellAutomationScript = provisionDir + "Setup.ps1"

func NewVMConstruct(winrmIP, winrmUsername, winrmPassword string) *VMConstruct {
	return &VMConstruct{NewWinRM(winrmIP, winrmUsername, winrmPassword)}
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
	fmt.Println("\nTransfering ~20 MB to the Windows VM. Depending on your connection, the transfer may take 15-45 minutes")
	err := c.uploadArtifact()
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

func (c *VMConstruct) uploadArtifact() error {
	fmt.Print("\tUploading LGPO to target VM...")
	err := c.remoteManager.UploadArtifact("./LGPO.zip", lgpoDest)
	if err != nil {
		return err
	}
	fmt.Println(" Finished uploading LGPO.")

	fmt.Print("\tUploading stemcell preparation artifacts to target VM...")
	err = c.remoteManager.UploadArtifact("./StemcellAutomation.zip", stemcellAutomationDest)
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
