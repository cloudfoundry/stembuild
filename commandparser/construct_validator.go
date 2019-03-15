package commandparser

import (
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/config"
	"os"
)

type ConstructValidator struct{}

func (c *ConstructValidator) PopulatedArgs(args ...string) bool {
	for _, arg := range args {
		if arg == "" {
			return false
		}
	}
	return true
}

func (c *ConstructValidator) LGPOInDirectory() bool {
	_, err := os.Stat("LGPO.zip")
	if err != nil {
		return false
	}
	return true
}

func (c *ConstructValidator) ValidStemcellInfo(stemcellVersion string) bool {
	return config.IsValidStemcellVersion(stemcellVersion)
}
