package version

import (
	"fmt"
	"strings"
)

type VersionGetterModifier interface {
	Modify(*VersionGetter)
}

func NewVersionGetter(modifiers ...VersionGetterModifier) *VersionGetter {
	v := &VersionGetter{
		Version: Version,
	}

	for _, modifier := range modifiers {
		modifier.Modify(v)
	}

	return v
}

type VersionGetter struct {
	Version string
}

func (v *VersionGetter) GetVersion() string {
	stringArr := strings.Split(v.Version, ".")
	stringArr = stringArr[0:2]

	return strings.Join(stringArr, ".")
}

func (v *VersionGetter) GetVersionWithPatchNumber(patchNumber string) string {
	return fmt.Sprintf("%s.%s", v.GetVersion(), patchNumber)
}

func (v *VersionGetter) GetOs() string {
	stringArr := strings.Split(v.Version, ".")
	os := stringArr[0]

	if os == "1200" {
		return "2012R2"
	}

	return os
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

func GetOSVersionFromBuildNumber(build string) string {
	for _, version := range AllVersions {
		if version.Build == build {
			return version.Name
		}
	}
	return ""
}
