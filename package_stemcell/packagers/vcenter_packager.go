package packagers

import (
	"errors"
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/stembuild/filesystem"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/config"
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/iaas_clients"
)

type VCenterPackager struct {
	SourceConfig config.SourceConfig
	Client       iaas_clients.IaasClient
}

func (v VCenterPackager) Package() error {
	err := v.Client.PrepareVM(v.SourceConfig.VmInventoryPath)
	if err != nil {
		return errors.New("could not prepare the VM for export")
	}
	return nil
}

func (v VCenterPackager) ValidateFreeSpaceForPackage(fs filesystem.FileSystem) error {
	println(os.Stdout, "WARNING: Please make sure you have enough free disk space for export")
	return nil
}

func (v VCenterPackager) ValidateSourceParameters() error {
	err := v.Client.ValidateUrl()
	if err != nil {
		return errors.New("please provide a valid vCenter URL")
	}

	err = v.Client.ValidateCredentials()
	if err != nil {
		errMsg := fmt.Sprintf("please provide valid credentials for %s", v.SourceConfig.URL)
		return errors.New(errMsg)
	}
	err = v.Client.FindVM(v.SourceConfig.VmInventoryPath)
	if err != nil {
		errorMsg := "VM path is invalid\nPlease make sure to format your inventory path correctly using the 'vm' keyword. Example: /my-datacenter/vm/my-folder/my-vm-name"
		return errors.New(errorMsg)
	}
	return nil
}
