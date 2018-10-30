package stembuild

import (
	"os"
	"regexp"
)

func IsValidVMDK(vmdk string) (bool, error) {

	if vmdk == "" {
		return false, nil
	}
	fi, err := os.Stat(vmdk)
	if err != nil {
		return false, err
	}
	if !fi.Mode().IsRegular() {
		return false, nil
	}

	return true, nil
}

func IsValidOS(os string) bool {
	switch os {
	case "2012R2", "1709", "1803":
		return true
	default:
		return false
	}
}

func IsValidVersion(version string) bool {

	if version == "" {
		return false
	}

	patterns := []string{
		`^\d{1,}\.\d{1,}$`,
		`^\d{1,}\.\d{1,}-build\.\d{1,}$`,
		`^\d{1,}\.\d{1,}\.\d{1,}$`,
		`^\d{1,}\.\d{1,}\.\d{1,}-build\.\d{1,}$`,
	}

	for _, pattern := range patterns {
		if regexp.MustCompile(pattern).MatchString(version) {
			return true
		}
	}

	return false
}
