package version

import (
	"strings"
)

func GetVersions(mainVersion string) (string, string) {
	stringArr := strings.Split(mainVersion, ".")

	os := stringArr[0]
	switch os {
	case "1709":
		os = "2016"
	case "1200":
		os = "2012R2"
	}
	stemcellVersion := strings.Join(stringArr[0:2], ".")

	return os, stemcellVersion
}
