package construct

import (
	"github.com/cloudfoundry-incubator/stembuild/remotemanager"
)

func NewMockVMConstruct(rm remotemanager.RemoteManager, client IaasClient, vmInventoryPath, username, password string) *VMConstruct {
	return &VMConstruct{rm, client, vmInventoryPath, username, password}
}

func (c *VMConstruct) UploadArtifact() error {
	return c.uploadArtifact()
}

func (c *VMConstruct) ExtractArchive() error {
	return c.extractArchive()
}

func (c *VMConstruct) ExecuteSetupScript() error {
	return c.executeSetupScript()
}
