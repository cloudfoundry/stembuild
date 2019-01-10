package commandparser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/stembuild/filesystem"
)

func HasAtLeastFreeDiskSpace(minFreeSpace uint64, fs filesystem.FileSystem, path string) (bool, uint64, error) {
	freeSpace, err := fs.GetAvailableDiskSpace(path)
	if err != nil {
		return false, 0, err
	}
	return freeSpace >= minFreeSpace, minFreeSpace - freeSpace, nil
}

func ValidateFreeSpaceForPackage(vmdkPath string, fs filesystem.FileSystem) (bool, uint64, error) {
	fi, err := os.Stat(vmdkPath)
	if err != nil {
		return false, uint64(0), fmt.Errorf("could not get vmdk info: %s", err)
	}
	vmdkSize := fi.Size()

	// make sure there is enough space for ova + stemcell and some leftover
	//	ova and stemcell will be the size of the vmdk in the worst case scenario
	minSpace := uint64(vmdkSize)*2 + (gigabyte / 2)
	hasSpace, spaceNeeded, err := HasAtLeastFreeDiskSpace(minSpace, fs, filepath.Dir(vmdkPath))
	if err != nil {
		return false, uint64(0), fmt.Errorf("could not check free space on disk: %s", err)
	}
	return hasSpace, spaceNeeded, nil
}
