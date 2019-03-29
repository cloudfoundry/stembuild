package construct

import (
	"github.com/cloudfoundry-incubator/stembuild/remotemanager"
)

func NewMockVMConstruct(
	rm remotemanager.RemoteManager,
	client IaasClient,
	vmInventoryPath,
	username,
	password string,
	unarchiver zipUnarchiver,
	messenger ConstructMessenger,
) *VMConstruct {

	return &VMConstruct{
		rm,
		client,
		vmInventoryPath,
		username,
		password,
		unarchiver,
		messenger,
	}
}
