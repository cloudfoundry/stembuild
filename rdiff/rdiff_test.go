package rdiff

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func createFile(name string) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func fileExists(name string) bool {
	fi, err := os.Stat(name)
	return err == nil && fi.Mode().IsRegular()
}

func TestPatchErrorHandling(t *testing.T) {
	dirname, err := ioutil.TempDir("", "rdiff-tests-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirname)

	basis := filepath.Join(dirname, "basis.txt")
	if err := createFile(basis); err != nil {
		t.Fatal(err)
	}

	delta := filepath.Join(dirname, "delta.txt")
	if err := createFile(delta); err != nil {
		t.Fatal(err)
	}

	doesNotExist := filepath.Join(dirname, "does-not-exist.txt")
	newfile := filepath.Join(dirname, "newfile.txt")

	// Make sure that an error is returned if the 'basis' or 'delta' files do
	// not exist and that newfile is not created if there is an error.
	//
	if err := Patch(doesNotExist, delta, newfile); err == nil {
		t.Error("Patch: expected error when the 'basis' file does not exist")
	}
	if fileExists(newfile) {
		t.Error("Patch: created newfile when the function arguments were invalid")
	}

	if err := Patch(basis, doesNotExist, newfile); err == nil {
		t.Error("Patch: expected error when the 'delta' file does not exist")
	}
	if fileExists(newfile) {
		t.Error("Patch: created newfile when the function arguments were invalid")
	}

	// Patch should return an error if newfile exists
	//
	if err := createFile(newfile); err != nil {
		t.Fatal(err)
	}
	switch err := Patch(basis, delta, newfile); err.(type) {
	case *os.PathError:
		// Ok
	default:
		t.Error("Patch: expected error when the 'newfile' file exists")
	}

	// Test that an error is returned from rsync, the files are empty
	// so should error.
	os.Remove(newfile)
	if err := Patch(basis, delta, newfile); err == nil {
		t.Error("Patch: expected error when the 'basis' and 'delta' files are invalid")
	}
}
