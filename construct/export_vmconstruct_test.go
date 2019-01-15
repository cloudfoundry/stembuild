package construct

import "github.com/cloudfoundry-incubator/stembuild/remotemanager"

func NewMockVMConstruct(rm remotemanager.RemoteManager) *VMConstruct {
	return &VMConstruct{rm}
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
