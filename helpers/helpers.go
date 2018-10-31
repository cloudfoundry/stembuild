package helpers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func recursiveFileList(destDir, searchDir string) ([]string, []string, []string, error) {
	srcFileList := make([]string, 0)
	destFileList := make([]string, 0)
	dirList := make([]string, 0)
	leafSearchDir := searchDir
	lastSepIndex := strings.LastIndex(searchDir, string(filepath.Separator))
	if lastSepIndex >= 0 {
		leafSearchDir = searchDir[lastSepIndex:len(searchDir)]
	}

	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			dirList = append(dirList, filepath.Join(destDir, leafSearchDir, path[len(searchDir):len(path)]))
		} else {
			srcFileList = append(srcFileList, path)
			destFileList = append(destFileList, filepath.Join(destDir, leafSearchDir, path[len(searchDir):len(path)]))
		}
		return err
	})

	if e != nil {
		return nil, nil, nil, e
	}

	return destFileList, srcFileList, dirList, nil
}

func CopyRecursive(destRoot, srcRoot string) error {
	var err error
	destRoot, err = filepath.Abs(destRoot)
	if err != nil {
		return err
	}

	srcRoot, err = filepath.Abs(srcRoot)
	if err != nil {
		return err
	}

	destFileList, srcFileList, dirList, err := recursiveFileList(destRoot, srcRoot)
	if err != nil {
		return err
	}

	// create destination directory hierarchy
	for _, myDir := range dirList {
		if err = os.MkdirAll(myDir, os.ModePerm); err != nil {
			return err
		}
	}

	for i, _ := range srcFileList {
		srcFile, err := os.Open(srcFileList[i])
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(destFileList[i])
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}

		if err = destFile.Sync(); err != nil {
			return err
		}
	}

	return nil
}

func CompareFiles(file1 string, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		return false
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false
	}
	defer f2.Close()

	blockSize := 64000
	block1 := make([]byte, blockSize)
	block2 := make([]byte, blockSize)

	for {
		_, err1 := f1.Read(block1)
		_, err2 := f2.Read(block2)

		switch {
		case err1 == io.EOF && err2 == io.EOF:
			return true
		case err1 != nil || err2 != nil:
			return false
		case !bytes.Equal(block1, block2):
			return false
		}
	}
}

func ExtractArchive(archive io.Reader, dirname string) error {
	tr := tar.NewReader(archive)

	limit := 100
	for ; limit >= 0; limit-- {
		h, err := tr.Next()
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("tar: reading from archive: %s", err)
			}
			break
		}

		// expect a flat archive
		name := h.Name
		if filepath.Base(name) != name {
			return fmt.Errorf("tar: archive contains subdirectory: %s", name)
		}

		// only allow regular files
		mode := h.FileInfo().Mode()
		if !mode.IsRegular() {
			return fmt.Errorf("tar: unexpected file mode (%s): %s", name, mode)
		}

		path := filepath.Join(dirname, name)
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
		if err != nil {
			return fmt.Errorf("tar: opening file (%s): %s", path, err)
		}
		defer f.Close()

		if _, err := io.Copy(f, tr); err != nil {
			return fmt.Errorf("tar: writing file (%s): %s", path, err)
		}
	}
	if limit <= 0 {
		return errors.New("tar: too many files in archive")
	}
	return nil
}

// ExtractGzipArchive extracts the tgz archive name to a temp directory
// returning the filepath of the temp directory.
func ExtractGzipArchive(name string) (string, error) {
	tmpdir, err := ioutil.TempDir("", "test-")
	if err != nil {
		return "", err
	}

	f, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	if err := ExtractArchive(w, tmpdir); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return tmpdir, nil
}

func ReadFile(name string) (string, error) {
	b, err := ioutil.ReadFile(name)
	return string(b), err
}
