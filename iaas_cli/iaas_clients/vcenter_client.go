package iaas_clients

import (
	"encoding/json"
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

func (c *VcenterClient) ValidateUrl() error {
	errCode := c.Runner.Run([]string{"about", "-u", c.Url})
	if errCode != 0 {
		return errors.New(fmt.Sprintf("vcenter_client - unable to validate url: %s", c.Url))
	}
	return nil

}

func (c *VcenterClient) ValidateCredentials() error {
	errCode := c.Runner.Run([]string{"about", "-u", c.credentialUrl})
	if errCode != 0 {
		return errors.New(fmt.Sprintf("vcenter_client - invalid credentials for: %s", c.Url))
	}

	return nil
}

func (c *VcenterClient) FindVM(vmInventoryPath string) error {
	errCode := c.Runner.Run([]string{"find", "-maxdepth=0", "-u", c.credentialUrl, vmInventoryPath})
	if errCode != 0 {
		return errors.New(fmt.Sprintf("vcenter_client - unable to find VM: %s. Ensure your inventory path is formatted properly and includes \"vm\" in its path, example: /my-datacenter/vm/my-folder/my-vm-name", vmInventoryPath))
	}

	return nil
}

func (c *VcenterClient) ListDevices(vmInventoryPath string) ([]string, error) {
	o, exitCode, err := c.Runner.RunWithOutput([]string{"device.ls", "-u", c.credentialUrl, "-vm", vmInventoryPath})

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
func (c *VcenterClient) RemoveDevice(vmInventoryPath string, deviceName string) error {
	errCode := c.Runner.Run([]string{"device.remove", "-u", c.credentialUrl, "-vm", vmInventoryPath, deviceName})
	if errCode != 0 {
		return fmt.Errorf("vcenter_client - %s could not be removed", deviceName)
	}
	return nil
}

func (c *VcenterClient) EjectCDRom(vmInventoryPath string, deviceName string) error {

	errCode := c.Runner.Run([]string{"device.cdrom.eject", "-u", c.credentialUrl, "-vm", vmInventoryPath, "-device", deviceName})
	if errCode != 0 {
		return fmt.Errorf("vcenter_client - %s could not be ejected", deviceName)
	}
	return nil
}

func (c *VcenterClient) ExportVM(vmInventoryPath string, destination string) error {
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

func (c *VcenterClient) UploadArtifact(vmInventoryPath, artifact, destination, username, password string) error {
	vmCredentials := fmt.Sprintf("%s:%s", username, password)
	errCode := c.Runner.Run([]string{"guest.upload", "-f", "-u", c.credentialUrl, "-l", vmCredentials, "-vm", vmInventoryPath, artifact, destination})
	if errCode != 0 {
		return fmt.Errorf("vcenter_client - %s could not be uploaded", artifact)
	}
	return nil
}

func (c *VcenterClient) MakeDirectory(vmInventoryPath, path, username, password string) error {
	vmCredentials := fmt.Sprintf("%s:%s", username, password)
	errCode := c.Runner.Run([]string{"guest.mkdir", "-u", c.credentialUrl, "-l", vmCredentials, "-vm", vmInventoryPath, "-p", path})
	if errCode != 0 {
		return fmt.Errorf("vcenter_client - directory `%s` could not be created", path)
	}
	return nil
}

func (c *VcenterClient) Start(vmInventoryPath, username, password, command string, args ...string) (string, error) {
	vmCredentials := fmt.Sprintf("%s:%s", username, password)
	pid, exitCode, err := c.Runner.RunWithOutput(append([]string{"guest.start", "-u", c.credentialUrl, "-l", vmCredentials, "-vm", vmInventoryPath, command}, args...))
	if err != nil {
		return "", fmt.Errorf("vcenter_client - failed to run '%s': %s", command, err)
	}
	if exitCode != 0 {
		return "", fmt.Errorf("vcenter_client - '%s' returned exit code: %d", command, exitCode)
	}
	// We trim this suffix since govc outputs the pid with an '\n' in the output
	return strings.TrimSuffix(pid, "\n"), nil
}

type govcPS struct {
	ProcessInfo []struct {
		Name      string
		Pid       int
		Owner     string
		CmdLine   string
		StartTime string
		EndTime   string
		ExitCode  int
	}
}

func (c *VcenterClient) WaitForExit(vmInventoryPath, username, password, pid string) (int, error) {
	vmCredentials := fmt.Sprintf("%s:%s", username, password)
	output, exitCode, err := c.Runner.RunWithOutput([]string{"guest.ps", "-u", c.credentialUrl, "-l", vmCredentials, "-vm", vmInventoryPath, "-p", pid, "-X", "-json"})
	if err != nil {
		return 0, fmt.Errorf("vcenter_client - failed to fetch exit code for PID %s: %s", pid, err)
	}
	if exitCode != 0 {
		return 0, fmt.Errorf("vcenter_client - fetching PID %s returned with exit code: %d", pid, exitCode)
	}

	ps := govcPS{}
	err = json.Unmarshal([]byte(output), &ps)
	if err != nil {
		return 0, fmt.Errorf("vcenter_client - received bad JSON output for PID %s: %s", pid, output)
	}
	if len(ps.ProcessInfo) != 1 {
		return 0, fmt.Errorf("vcenter_client - couldn't get exit code for PID %s", pid)
	}

	return ps.ProcessInfo[0].ExitCode, nil
}
