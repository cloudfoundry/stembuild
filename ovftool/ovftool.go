package ovftool

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var stopWalk = errors.New("stop walk")

func findExecutable(root, name string) (string, error) {
	var file string
	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() && fi.Name() == name {
			if s, err := exec.LookPath(path); err == nil {
				file = s
				return stopWalk
			}
		}
		return nil
	}
	err := filepath.Walk(root, walkFn)
	if file == "" {
		if err == nil || err == stopWalk {
			err = fmt.Errorf("executable file not found in: %s", root)
		}
		// CEV: this should never happen
		if err == stopWalk {
			err = fmt.Errorf("executable file not found in: %s - exec.LookPath error.", root)
		}
		return "", &exec.Error{Name: name, Err: err}
	}
	return file, nil
}
