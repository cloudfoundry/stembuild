package factory

import (
	"errors"
	"os"
	"strings"

	"github.com/cloudfoundry-incubator/stembuild/colorlogger"
	"github.com/cloudfoundry-incubator/stembuild/pack/config"
	"github.com/cloudfoundry-incubator/stembuild/pack/options"
	"github.com/cloudfoundry-incubator/stembuild/pack/stemcell"
)

type Packager interface {
	Package() error
	ValidateSourceParameters() error
}

func GetPackager(sourceConfig config.SourceConfig, outputConfig config.OutputConfig, logLevel int, color bool) (Packager, error) {
	source, err := sourceConfig.GetSource()
	if err != nil {
		return nil, err
	}
	switch source {
	case config.VCENTER:
		v := VCenterPackager{SourceConfig: sourceConfig}
		return v, nil
	case config.VMDK:
		options := options.StembuildOptions{}
		logger := colorlogger.ConstructLogger(logLevel, color, os.Stderr)
		vmdkPackager := stemcell.Config{
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
