package version

import (
	"strings"
)

func GetVersions(mainVersion string) (string, string) {
	stringArr := strings.Split(mainVersion, ".")

	os := stringArr[0]
	if os == "1709" {
		os = "2016"
	}
	stemcellVersion := strings.Join(stringArr[0:2], ".")

	return os, stemcellVersion
}