package iaas_clients

import (
	"errors"
	"fmt"

	"github.com/cloudfoundry-incubator/stembuild/iaas_cli"
)

//go:generate counterfeiter . IaasClient
type IaasClient interface {
	ValidateUrl() error
	ValidateCredentials() error
	FindVM(vmInventoryPath string) error
	PrepareVM(vmInventoryPath string) error
}

type VcenterClient struct {
	Username      string
	Password      string
	Url           string
	credentialUrl string
	Runner        iaas_cli.CliRunner
}

func NewVcenterClient(username string, password string, url string, runner iaas_cli.CliRunner) *VcenterClient {
	urlWithCredentials := fmt.Sprintf("%s:%s@%s", username, password, url)
	return &VcenterClient{Username: username, Password: password, Url: url, credentialUrl: urlWithCredentials, Runner: runner}
}

func (c VcenterClient) ValidateUrl() error {
	errCode := c.Runner.Run([]string{"about", "-u", c.Url})
	if errCode != 0 {
		return errors.New("invalid url")
	}
	return nil

}

func (c VcenterClient) ValidateCredentials() error {
	errCode := c.Runner.Run([]string{"about", "-u", c.credentialUrl})
	if errCode != 0 {
		return errors.New("invalid credentials")
	}

	return nil
}

func (c VcenterClient) FindVM(vmInventoryPath string) error {
	errCode := c.Runner.Run([]string{"find", "-maxdepth=0", "-u", c.credentialUrl, vmInventoryPath})
	if errCode != 0 {
		errorMsg := "invalid VM path"
		return errors.New(errorMsg)
	}

	return nil
}

func (c VcenterClient) PrepareVM(vmInventoryPath string) error {
	ethernetDeviceName, floppyDeviceName := "ethernet-0", "floppy-8000"
	var err error

	err = c.removeDevice(vmInventoryPath, ethernetDeviceName)
	if err != nil {
		return err
	}

	err = c.removeDevice(vmInventoryPath, floppyDeviceName)
	if err != nil {
		return err
	}

	return nil
}

func (c VcenterClient) removeDevice(vmInventoryPath string, deviceName string) error {
	errCode := c.Runner.Run([]string{"device.remove", "-vm", vmInventoryPath, deviceName, "-u", c.credentialUrl})
	if errCode != 0 {
		errorMsg := fmt.Sprintf(deviceName + " could not be removed/not found")
		return errors.New(errorMsg)
	}
	return nil
}
