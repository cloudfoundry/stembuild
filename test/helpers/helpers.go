package helpers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

func extractArchive(archive io.Reader, dirname string) error {
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
	if err := extractArchive(w, tmpdir); err != nil {
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

func BuildStembuild() (string, error) {
	stembuildExecutable, err := ioutil.TempFile(os.TempDir(), "stembuild")
	if err != nil {
		return "", err
	}

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})

	buildCommand := fmt.Sprintf("go build -o %s %s", stembuildExecutable.Name(), "github.com/cloudfoundry-incubator/stembuild")
	buildCommandSlice := strings.Split(buildCommand, " ")

	cmd := exec.Command(buildCommandSlice[0], buildCommandSlice[1:]...)
	session, err := gexec.Start(cmd, stdout, stderr)
	if err != nil {
		return "", err
	}
	gomega.EventuallyWithOffset(1, session, 30*time.Second).Should(
		gexec.Exit(0),
		fmt.Sprintf(
			"Build command %s exited with exit code: %d, stdout: %s, stderr: %s",
			buildCommand,
			session.ExitCode(),
			string(stdout.Bytes()),
			string(stderr.Bytes()),
		),
	)

	return stembuildExecutable.Name(), err
}
