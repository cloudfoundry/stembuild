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

func (p VCenterPackager) ValidateFreeSpaceForPackage(fs filesystem.FileSystem) error {
	return nil
}

func (v VCenterPackager) ValidateSourceParameters() error {
	err := v.Client.ValidateUrl()
	if err != nil {
		return errors.New("please provide a valid vCenter URL")
	}

	err = v.Client.Login()
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

type VcenterClient interface {
	ValidateUrl() error
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

func (c FakeVcenterClient) ValidateUrl() error {
	if c.InvalidUrl {
		return errors.New("invalid url")
	}
	return nil
}

func (c FakeVcenterClient) Login() error {
	if c.InvalidCredentials {
		return errors.New("invalid credentials")
	}
	return nil
}

func (c FakeVcenterClient) FindVM(vmInventoryPath string) error {
	if c.InvalidVmInventoryPath {
		return errors.New("invalid VM path")
	}
	return nil
}

type RealVcenterClient struct {
	Username      string
	Password      string
	Url           string
	credentialUrl string
}

func NewRealVcenterClient(username string, password string, url string) *RealVcenterClient {
	urlWithCredentials := fmt.Sprintf("%s:%s@%s", username, password, url)
	return &RealVcenterClient{Username: username, Password: password, Url: url, credentialUrl: urlWithCredentials}
}

func (c RealVcenterClient) ValidateUrl() error {
	errCode := cli.Run([]string{"about", "-u", c.Url})
	if errCode != 0 {
		return errors.New("invalid url")
	}
	return nil

}

func (c RealVcenterClient) Login() error {
	errCode := cli.Run([]string{"about", "-u", c.credentialUrl})
	if errCode != 0 {
		return errors.New("invalid credentials")
	}

	return nil
}

func (c RealVcenterClient) FindVM(vmInventoryPath string) error {
	errCode := cli.Run([]string{"find", "-maxdepth=0", "-u", c.credentialUrl, vmInventoryPath})
	if errCode != 0 {
		errorMsg := "invalid VM path"
		return errors.New(errorMsg)
	}

	return nil
}
