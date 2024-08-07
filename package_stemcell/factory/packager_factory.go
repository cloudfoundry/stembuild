package factory

import (
	"errors"
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

func (f *PackagerFactory) Packager(sourceConfig config.SourceConfig, outputConfig config.OutputConfig, logger colorlogger.Logger) (commandparser.Packager, error) {
	source, err := sourceConfig.GetSource()
	if err != nil {
		return nil, err
	}

	switch source {
	case config.VCENTER:
		client :=
			iaas_clients.NewVcenterClient(
				sourceConfig.Username,
				sourceConfig.Password,
				sourceConfig.URL,
				sourceConfig.CaCertFile,
				&iaas_cli.GovcRunner{},
			)

		return &packagers.VCenterPackager{
			SourceConfig: sourceConfig,
			OutputConfig: outputConfig,
			Client:       client,
			Logger:       logger,
		}, nil
	case config.VMDK:
		options :=
			package_parameters.VmdkPackageParameters{}

		vmdkPackager := &packagers.VmdkPackager{
			Stop:         make(chan struct{}),
			BuildOptions: options,
			Logger:       logger,
		}

		vmdkPackager.BuildOptions.VMDKFile = sourceConfig.Vmdk
		vmdkPackager.BuildOptions.OSVersion = strings.ToUpper(outputConfig.Os)
		vmdkPackager.BuildOptions.Version = outputConfig.StemcellVersion
		vmdkPackager.BuildOptions.OutputDir = outputConfig.OutputDir
		return vmdkPackager, nil
	default:
		return nil, errors.New("unable to determine packager")
	}
}
