package iaas_clients

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cloudfoundry-incubator/stembuild/iaas_cli"
)

type VcenterClient struct {
	Url           string
	credentialUrl string
	Runner        iaas_cli.CliRunner
}

func NewVcenterClient(username string, password string, url string, runner iaas_cli.CliRunner) *VcenterClient {
	urlWithCredentials := fmt.Sprintf("%s:%s@%s", username, password, url)
	return &VcenterClient{Url: url, credentialUrl: urlWithCredentials, Runner: runner}
}

func (c VcenterClient) ValidateUrl() error {
	errCode := c.Runner.Run([]string{"about", "-u", c.Url})
	if errCode != 0 {
		return errors.New(fmt.Sprintf("vcenter_client - unable to validate url: %s", c.Url))
	}
	return nil

}

func (c VcenterClient) ValidateCredentials() error {
	errCode := c.Runner.Run([]string{"about", "-u", c.credentialUrl})
	if errCode != 0 {
		return errors.New(fmt.Sprintf("vcenter_client - invalid credentials for: %s", c.Url))
	}

	return nil
}

func (c VcenterClient) FindVM(vmInventoryPath string) error {
	errCode := c.Runner.Run([]string{"find", "-maxdepth=0", "-u", c.credentialUrl, vmInventoryPath})
	if errCode != 0 {
		return errors.New(fmt.Sprintf("vcenter_client - unable to find VM: %s. Ensure your inventory path is formatted properly and includes \"vm\" in its path, example: /my-datacenter/vm/my-folder/my-vm-name", vmInventoryPath))
	}

	return nil
}

func (c VcenterClient) ListDevices(vmInventoryPath string) ([]string, error) {
	o, exitCode, err := c.Runner.RunWithOutput([]string{"device.ls", "-vm", vmInventoryPath})

	if exitCode != 0 {
		return []string{}, fmt.Errorf("vcenter_client - failed to list devices in vCenter, govc exit code %d", exitCode)
	}

	if err != nil {
		return []string{}, fmt.Errorf("vcenter_client - failed to parse list of devices. Err: %s", err)
	}

	entries := strings.Split(o, "\n")
	devices := []string{}
	r, _ := regexp.Compile(`\S+`)
	for _, entry := range entries {
		if entry != "" {
			devices = append(devices, r.FindString(entry))
		}
	}
	return devices, nil
}
func (c VcenterClient) RemoveDevice(vmInventoryPath string, deviceName string) error {
	errCode := c.Runner.Run([]string{"device.remove", "-u", c.credentialUrl, "-vm", vmInventoryPath, deviceName})
	if errCode != 0 {
		return fmt.Errorf("vcenter_client - %s could not be removed", deviceName)
	}
	return nil
}

func (c VcenterClient) EjectCDRom(vmInventoryPath string, deviceName string) error {

	errCode := c.Runner.Run([]string{"device.cdrom.eject", "-u", c.credentialUrl, "-vm", vmInventoryPath, "-device", deviceName})
	if errCode != 0 {
		return fmt.Errorf("vcenter_client - %s could not be ejected", deviceName)
	}
	return nil
}

func (c VcenterClient) ExportVM(vmInventoryPath string, destination string) error {
	_, err := os.Stat(destination)
	if err != nil {
		return errors.New(fmt.Sprintf("vcenter_client - provided destination directory: %s does not exist", destination))
	}
	errCode := c.Runner.Run([]string{"export.ovf", "-u", c.credentialUrl, "-sha", "1", "-vm", vmInventoryPath, destination})
	if errCode != 0 {
		return errors.New(fmt.Sprintf("vcenter_client - %s could not be exported", vmInventoryPath))
	}
	return nil
}
