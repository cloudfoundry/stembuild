// +build !windows

package main

import "io/ioutil"

func TempDir(prefix string) (string, error) {
	return ioutil.TempDir("", prefix)
}
