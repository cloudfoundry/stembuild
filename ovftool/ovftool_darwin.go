// +build darwin

package ovftool

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
)

// homeDirectory returns the home directory of the current user,
// errors are ignored.
func homeDirectory() string {
	if s := os.Getenv("HOME"); s != "" {
		return s
	}

	out, err := exec.Command("sh", "-c", "cd ~ && pwd").Output()
	if err != nil {
		return ""
	}

	s := string(bytes.TrimSpace(out))
	if s == "" {
		return ""
	}
	return s
}

func Ovftool() (string, error) {
	const name = "ovftool"
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}

	// search paths
	var vmwareDirs = []string{
		"/Applications/VMware Fusion.app",
	}
	if home := homeDirectory(); home != "" {
		vmwareDirs = append(vmwareDirs, filepath.Join(home, vmwareDirs[0]))
	}

	for _, root := range vmwareDirs {
		if fi, err := os.Stat(root); err != nil || !fi.IsDir() {
			continue
		}
		if path, err := findExecutable(root, name); err == nil {
			return path, nil
		}
	}

	return "", &exec.Error{Name: name, Err: exec.ErrNotFound}
}
