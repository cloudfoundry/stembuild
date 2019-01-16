package packagers

import (
	"errors"
	"fmt"

	"github.com/cloudfoundry-incubator/stembuild/filesystem"

	_ "github.com/vmware/govmomi/govc/about"
	"github.com/vmware/govmomi/govc/cli"
	_ "github.com/vmware/govmomi/govc/object"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/config"
)

type VCenterPackager struct {
	SourceConfig config.SourceConfig
	Client       VcenterClient
}

func (v VCenterPackager) Package() error {
	return nil
}

func (v VCenterPackager) ValidateFreeSpaceForPackage(fs filesystem.FileSystem) error {
	return nil
}

func (v VCenterPackager) ValidateSourceParameters() error {
	err := v.Client.Login()
	if err != nil {
		return err
	}

	err = v.Client.FindVM(v.SourceConfig.VmInventoryPath)
	if err != nil {
		return err
	}
	return nil
}

type VcenterClient interface {
	Login() error
	FindVM(vmInventoryPath string) error
}

type FakeVcenterClient struct {
	Username               string
	Password               string
	Url                    string
	InvalidCredentials     bool
	InvalidUrl             bool
	InvalidVmInventoryPath bool
}

func (c FakeVcenterClient) Login() error {
	if c.InvalidCredentials {
		errorMsg := fmt.Sprintf("please provide valid credentials for %s", c.Url)
		return errors.New(errorMsg)
	}
	if c.InvalidUrl {
		return errors.New("please provide a valid vCenter URL")
	}
	return nil
}

func (c FakeVcenterClient) FindVM(vmInventoryPath string) error {
	if c.InvalidVmInventoryPath {
		return errors.New("VM path is invalid\nPlease make sure to format your inventory path correctly using the 'vm' keyword. Example: /my-datacenter/vm/my-folder/my-vm-name")
	}
	return nil
}

type RealVcenterClient struct {
	Username string
	Password string
	Url      string
}

func (c RealVcenterClient) Login() error {

	errCode := cli.Run([]string{"about", "-u", c.Url})
	if errCode != 0 {
		return errors.New("please provide valid vCenter URL")
	}

	errCode = cli.Run([]string{"about", "-u", fmt.Sprintf("%s:%s@%s", c.Username, c.Password, c.Url)})
	if errCode != 0 {
		errorMsg := fmt.Sprintf("please provide valid credentials for %s", c.Url)
		return errors.New(errorMsg)
	}

	return nil
}

func (c RealVcenterClient) FindVM(vmInventoryPath string) error {
	errCode := cli.Run([]string{"find", "-maxdepth=0", "-u", fmt.Sprintf("%s:%s@%s", c.Username, c.Password, c.Url), vmInventoryPath})
	if errCode != 0 {
		errorMsg := "VM path is invalid\nPlease make sure to format your inventory path correctly using the 'vm' keyword. Example: /my-datacenter/vm/my-folder/my-vm-name"
		return errors.New(errorMsg)
	}

	return nil
}
