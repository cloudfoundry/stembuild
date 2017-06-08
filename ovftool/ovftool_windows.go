package ovftool

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

// vmwareInstallPaths, returns the install paths of VMware Workstation and
// OVF Tool, which can be installed separately.
func vmwareInstallPaths() ([]string, error) {
	const regKey = registry.LOCAL_MACHINE
	const access = registry.QUERY_VALUE

	keypaths := []string{
		`SOFTWARE\Wow6432Node\VMware, Inc.\VMware Workstation`,
		`SOFTWARE\Wow6432Node\VMware, Inc.\VMware OVF Tool`,
		`SOFTWARE\VMware, Inc.\VMware Workstation`,
		`SOFTWARE\VMware, Inc.\VMware OVF Tool`,
	}

	var key registry.Key
	var err error
	for _, path := range keypaths {
		key, err = registry.OpenKey(regKey, path, access)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("opening VMware Workstation and OVF Tool registry keys: %s", err)
	}
	defer key.Close()

	var first error
	var paths []string
	for _, k := range []string{"InstallPath64", "InstallPath"} {
		s, _, err := key.GetStringValue(k)
		if err != nil && first == nil {
			first = err
		} else {
			paths = append(paths, s)
		}
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("could not find VMware Workstation install path in registry:", first)
	}
	return paths, nil
}

func Ovftool() (string, error) {
	const name = "ovftool.exe"
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}

	installPaths, err := vmwareInstallPaths()
	if err != nil {
		return "", err
	}

	// Locations of ovftool.exe in the OVF Tool and Workstation directories
	search := []string{
		// Location if OVF Tool is installed
		name,
		// Location if Workstation is installed
		filepath.Join("ovftool", name),
	}
	for _, dir := range installPaths {
		for _, name := range search {
			file := filepath.Join(dir, name)
			if path, err := exec.LookPath(file); err == nil {
				return path, nil
			}
		}
		if path, err := findExecutable(dir, name); err == nil {
			return path, nil
		}
	}
	return "", &exec.Error{Name: name, Err: exec.ErrNotFound}
}
