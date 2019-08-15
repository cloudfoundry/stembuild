package version

import (
	"strings"
)

type VersionGetter struct{}

func (v *VersionGetter) GetVersion() string {
	return Version
}

type WindowsVersion struct {
	Name  string
	Build string
}

var Version = "dev"
var (
	Version2019 = WindowsVersion{Name: "2019", Build: "17763"}
	Version1803 = WindowsVersion{Name: "1803", Build: "17134"}
	VersionDev  = WindowsVersion{Name: "dev", Build: "dev"}

	AllVersions = []WindowsVersion{Version2019, Version1803, VersionDev}
)

func GetVersions(mainVersion string) (string, string) {
	stringArr := strings.Split(mainVersion, ".")

	// TODO remove special-case handling when we stop building 2012 stemcells
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

func GetOSVersionFromBuildNumber(build string) string {
	for _, version := range AllVersions {
		if version.Build == build {
			return version.Name
		}
	}
	return ""
}
