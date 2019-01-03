package commandparser

import (
	"fmt"
	"github.com/cloudfoundry-incubator/stembuild/filesystem"
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
	case "2012R2", "1803", "2016":
		return true
	default:
		return false
	}
}

func ValidateOrCreateOutputDir(outputDir string) error {

	fi, err := os.Stat(outputDir)
	if err != nil && os.IsNotExist(err) {
		if err = os.Mkdir(outputDir, 0700); err != nil {
			return err
		}
	} else if err != nil || fi == nil {
		return fmt.Errorf("error opening output directory (%s): %s\n", outputDir, err)
	} else if !fi.IsDir() {
		return fmt.Errorf("output argument (%s): is not a directory\n", outputDir)
	}

	return nil
}

func IsValidStemcellVersion(version string) bool {

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

func HasAtLeastFreeDiskSpace(minFreeSpace uint64, fs filesystem.FileSystem, path string) (bool, uint64, error) {
	freeSpace, err := fs.GetAvailableDiskSpace(path)
	if err != nil {
		return false, 0, err
	}
	return freeSpace >= minFreeSpace, minFreeSpace - freeSpace, nil
}
