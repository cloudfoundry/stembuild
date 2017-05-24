package ovftool

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func workstationInstallPaths() ([]string, error) {
	const (
		keypath1 = `SOFTWARE\Wow6432Node\VMware, Inc.\VMware Workstation`
		keypath2 = `SOFTWARE\VMware, Inc.\VMware Workstation`
		regKey   = registry.LOCAL_MACHINE
		access   = registry.QUERY_VALUE
	)
	key, e1 := registry.OpenKey(regKey, keypath1, access)
	if e1 != nil {
		var e2 error
		key, e2 = registry.OpenKey(regKey, keypath2, access)
		if e2 != nil {
			return nil, fmt.Errorf("opening VMware Work Station registry keys: (%s): %s; (%s): %s",
				keypath1, e1, keypath2, e2)
		}
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
		return nil, fmt.Errorf("could not find VMware Work Station install path in registry:", first)
	}
	return paths, nil
}

func Ovftool() (string, error) {
	const name = "ovftool.exe"
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}

	installPaths, err := workstationInstallPaths()
	if err != nil {
		return "", err
	}
	for _, dir := range installPaths {
		file := filepath.Join(dir, "ovftool", "ovftool.exe")
		if path, err := exec.LookPath(file); err == nil {
			return path, nil
		}
		if path, err := findExecutable(dir, name); err == nil {
			return path, nil
		}
	}
	return "", &exec.Error{Name: name, Err: exec.ErrNotFound}
}
