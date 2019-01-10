package packagers

import (
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/config"
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
