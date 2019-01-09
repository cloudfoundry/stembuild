package factory

import (
	"github.com/cloudfoundry-incubator/stembuild/pack/config"
)

type VCenterPackager struct {
	SourceConfig config.SourceConfig
}

func (v VCenterPackager) Package() error {
	return nil
}

func (v VCenterPackager) ValidateSourceParameters() error {
	return nil
}
