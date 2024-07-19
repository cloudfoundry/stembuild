package factory

import (
	"errors"
	"os"
	"strings"

	"github.com/cloudfoundry/stembuild/colorlogger"
	"github.com/cloudfoundry/stembuild/commandparser"
	"github.com/cloudfoundry/stembuild/iaas_cli"
	"github.com/cloudfoundry/stembuild/iaas_cli/iaas_clients"
	"github.com/cloudfoundry/stembuild/package_stemcell/config"
	"github.com/cloudfoundry/stembuild/package_stemcell/package_parameters"
	"github.com/cloudfoundry/stembuild/package_stemcell/packagers"
)

type PackagerFactory struct{}

func (f *PackagerFactory) Packager(sourceConfig config.SourceConfig, outputConfig config.OutputConfig, logLevel int, color bool) (commandparser.Packager, error) {
	source, err := sourceConfig.GetSource()
	if err != nil {
		return nil, err
	}
	switch source {
	case config.VCENTER:
		runner := &iaas_cli.GovcRunner{}
		client := iaas_clients.NewVcenterClient(sourceConfig.Username, sourceConfig.Password, sourceConfig.URL, sourceConfig.CaCertFile, runner)
		vCenterPackager := &packagers.VCenterPackager{SourceConfig: sourceConfig, OutputConfig: outputConfig, Client: client}
		return vCenterPackager, nil
	case config.VMDK:
		options := package_parameters.VmdkPackageParameters{}
		logger := colorlogger.ConstructLogger(logLevel, color, os.Stderr)
		vmdkPackager := &packagers.VmdkPackager{
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
