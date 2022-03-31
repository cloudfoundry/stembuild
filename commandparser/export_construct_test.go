package commandparser

import "github.com/cloudfoundry/stembuild/construct/config"

func (p *ConstructCmd) GetSourceConfig() config.SourceConfig {
	return p.sourceConfig
}
