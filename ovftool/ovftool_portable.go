// +build !darwin,!windows

package ovftool

import "os/exec"

func Ovftool() (string, error) {
	const name = "ovftool"
	return exec.LookPath(name)
}
