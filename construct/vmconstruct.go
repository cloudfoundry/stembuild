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

	fmt.Println("validating connection to vm...")
	err := c.remoteManager.CanReachVM()
	if err != nil {
		return err
	}
	return c.remoteManager.CanLoginVM()
}

func (c *VMConstruct) PrepareVM() error {

	fmt.Println("upload artifact...")
	err := c.uploadArtifact()
	if err != nil {
		return err
	}
	fmt.Println("extract artifact...")
	err = c.extractArchive()
	if err != nil {
		return err
	}
	fmt.Println("execute script...")
	err = c.executeSetupScript()
	if err != nil {
		return err
	}

	return nil
}

func (c *VMConstruct) uploadArtifact() error {
	err := c.remoteManager.UploadArtifact("./LGPO.zip", lgpoDest)
	if err != nil {
		return err
	}
	err = c.remoteManager.UploadArtifact("./StemcellAutomation.zip", stemcellAutomationDest)
	if err != nil {
		return err
	}

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
