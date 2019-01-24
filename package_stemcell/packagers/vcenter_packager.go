package packagers

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/cloudfoundry-incubator/stembuild/filesystem"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/config"
)

//go:generate counterfeiter . IaasClient
type IaasClient interface {
	ValidateUrl() error
	ValidateCredentials() error
	FindVM(vmInventoryPath string) error
	ExportVM(vmInventoryPath string, destination string) error
	ListDevices(vmInventoryPath string) ([]string, error)
	RemoveDevice(vmInventoryPath string, deviceName string) error
	EjectCDRom(vmInventoryPath string, deviceName string) error
}

type VCenterPackager struct {
	SourceConfig config.SourceConfig
	OutputConfig config.OutputConfig
	Client       IaasClient
}

func (v VCenterPackager) Package() error {
	err := v.executeOnMatchingDevice(v.Client.RemoveDevice, "^(floppy-|ethernet-)")
	if err != nil {
		return err
	}
	err = v.executeOnMatchingDevice(v.Client.EjectCDRom, "^(cdrom-)")
	if err != nil {
		return err
	}

	workingDir, err := ioutil.TempDir(os.TempDir(), "vcenter-packager-working-directory")

	if err != nil {
		return errors.New("failed to create working directory")
	}

	stemcellDir, err := ioutil.TempDir(os.TempDir(), "vcenter-packager-stemcell-directory")
	if err != nil {
		return errors.New("failed to create stemcell directory")
	}
	err = v.Client.ExportVM(v.SourceConfig.VmInventoryPath, workingDir)

	if err != nil {
		return errors.New("failed to export the prepared VM")
	}

	vmName := path.Base(v.SourceConfig.VmInventoryPath)
	shaSum, err := TarGenerator(filepath.Join(stemcellDir, "image"), filepath.Join(workingDir, vmName))
	manifestContents := CreateManifest(v.OutputConfig.Os, v.OutputConfig.StemcellVersion, shaSum)
	err = WriteManifest(manifestContents, stemcellDir)

	if err != nil {
		return errors.New("failed to create stemcell.MF file")
	}

	stemcellFilename := StemcellFilename(v.OutputConfig.StemcellVersion, v.OutputConfig.Os)
	_, err = TarGenerator(filepath.Join(v.OutputConfig.OutputDir, stemcellFilename), stemcellDir)

	return nil
}

func (v VCenterPackager) executeOnMatchingDevice(action func(a, b string) error, devicePattern string) error {
	deviceList, err := v.Client.ListDevices(v.SourceConfig.VmInventoryPath)
	if err != nil {
		return err
	}

	for _, deviceName := range deviceList {
		matched, _ := regexp.MatchString(devicePattern, deviceName)
		if matched {
			err = action(v.SourceConfig.VmInventoryPath, deviceName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (v VCenterPackager) ValidateFreeSpaceForPackage(fs filesystem.FileSystem) error {
	return nil
}

func (v VCenterPackager) ValidateSourceParameters() error {
	err := v.Client.ValidateUrl()
	if err != nil {
		return err
	}

	err = v.Client.ValidateCredentials()
	if err != nil {
		return err
	}
	err = v.Client.FindVM(v.SourceConfig.VmInventoryPath)
	if err != nil {
		return err
	}
	return nil
}
