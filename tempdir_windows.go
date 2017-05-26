package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

func TempDir(prefix string) (string, error) {
	windir := os.Getenv("WINDIR")
	if windir == "" {
		return "", errors.New("missing WINDIR environment variable")
	}
	dir := filepath.Join(windir, "Temp")
	return ioutil.TempDir(dir, prefix)
}
