package commandparser

import "github.com/cloudfoundry-incubator/stembuild/construct/config"

func (p *ConstructCmd) GetSourceConfig() config.SourceConfig {
	return p.sourceConfig
}
