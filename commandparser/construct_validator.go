package commandparser

import (
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/config"
	"os"
	"path/filepath"
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
	dir, _ := os.Getwd()
	_, err := os.Stat(filepath.Join(dir, "LGPO.zip"))
	if err != nil {
		return false
	}
	return true
}

func (c *ConstructValidator) ValidStemcellInfo(stemcellVersion string) bool {
	return config.IsValidStemcellVersion(stemcellVersion)
}
