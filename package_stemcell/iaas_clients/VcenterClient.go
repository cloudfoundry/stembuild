package iaas_clients

import (
	"errors"
	"fmt"

	_ "github.com/vmware/govmomi/govc/about"
	"github.com/vmware/govmomi/govc/cli"
	_ "github.com/vmware/govmomi/govc/object"
)

//go:generate counterfeiter . VcenterClient
type VcenterClient interface {
	ValidateUrl() error
	Login() error
	FindVM(vmInventoryPath string) error
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
