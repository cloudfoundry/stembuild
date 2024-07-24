package packager

import (
	"errors"
	"strings"

	"github.com/cloudfoundry/stembuild/colorlogger"
	"github.com/cloudfoundry/stembuild/commandparser"
	"github.com/cloudfoundry/stembuild/iaas_cli"
	"github.com/cloudfoundry/stembuild/iaas_cli/iaas_clients"
	"github.com/cloudfoundry/stembuild/package_stemcell/config"
	"github.com/cloudfoundry/stembuild/package_stemcell/package_parameters"
)

type Factory struct{}

func (f *Factory) NewPackager(sourceConfig config.SourceConfig, outputConfig config.OutputConfig, logger colorlogger.Logger) (commandparser.Packager, error) {
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

		return &VCenterPackager{
			SourceConfig: sourceConfig,
			OutputConfig: outputConfig,
			Client:       client,
			Logger:       logger,
		}, nil
	case config.VMDK:
		options :=
			package_parameters.VmdkPackageParameters{
				VMDKFile:  sourceConfig.Vmdk,
				OSVersion: strings.ToUpper(outputConfig.Os),
				Version:   outputConfig.StemcellVersion,
				OutputDir: outputConfig.OutputDir,
			}

		return &VmdkPackager{
			Stop:         make(chan struct{}),
			BuildOptions: options,
			Logger:       logger,
		}, nil
	default:
		return nil, errors.New("unable to determine packager")
	}
}
