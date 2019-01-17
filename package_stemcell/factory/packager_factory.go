package factory

import (
	"errors"
	"os"
	"strings"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/iaas_clients"

	"github.com/cloudfoundry-incubator/stembuild/filesystem"

	"github.com/cloudfoundry-incubator/stembuild/colorlogger"
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/config"
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/package_parameters"
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/packagers"
)

type Packager interface {
	Package() error
	ValidateFreeSpaceForPackage(fs filesystem.FileSystem) error
	ValidateSourceParameters() error
}

func GetPackager(sourceConfig config.SourceConfig, outputConfig config.OutputConfig, logLevel int, color bool) (Packager, error) {
	source, err := sourceConfig.GetSource()
	if err != nil {
		return nil, err
	}
	switch source {
	case config.VCENTER:
		client := iaas_clients.NewRealVcenterClient(sourceConfig.Username, sourceConfig.Password, sourceConfig.URL)
		v := packagers.VCenterPackager{SourceConfig: sourceConfig, Client: client}
		return v, nil
	case config.VMDK:
		options := package_parameters.VmdkPackageParameters{}
		logger := colorlogger.ConstructLogger(logLevel, color, os.Stderr)
		vmdkPackager := packagers.VmdkPackager{
			Stop:         make(chan struct{}),
			Debugf:       logger.Debugf,
			BuildOptions: options,
		}

		vmdkPackager.BuildOptions.VMDKFile = sourceConfig.Vmdk
		vmdkPackager.BuildOptions.OSVersion = strings.ToUpper(outputConfig.Os)
		vmdkPackager.BuildOptions.Version = outputConfig.StemcellVersion
		vmdkPackager.BuildOptions.OutputDir = outputConfig.OutputDir
		return vmdkPackager, nil
	}
	return nil, errors.New("Unable to determine packager")
}
