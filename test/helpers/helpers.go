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
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/onsi/ginkgo"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
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

func BuildStembuild(version string) (string, error) {
	var stembuildExe = "stembuild"
	if runtime.GOOS == "windows" {
		stembuildExe += ".exe"
	}

	tempDir, err := ioutil.TempDir("", "stembuild")
	if err != nil {
		return "", err
	}

	stembuildExecutable := fmt.Sprintf("%s%c%s", tempDir, os.PathSeparator, stembuildExe)

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})

	buildPackage := "github.com/cloudfoundry-incubator/stembuild"
	buildPackage = filepath.FromSlash(buildPackage)

	generateCommand := fmt.Sprintf("go generate %s", buildPackage)
	generateCommandSlice := strings.Split(generateCommand, " ")

	genCmd := exec.Command(generateCommandSlice[0], generateCommandSlice[1:]...)
	genSess, err := gexec.Start(genCmd, stdout, stderr)
	if err != nil {
		return "", err
	}
	gomega.EventuallyWithOffset(1, genSess, 30*time.Second).Should(
		gexec.Exit(0),
		fmt.Sprintf("Generate command failed with exit code %d", genSess.ExitCode()),
	)

	goPath := EnvMustExist("GOPATH")
	prefix := filepath.Join(goPath, "src", buildPackage, "integration", "construct", "assets")
	stemcellAutomation := filepath.Join(prefix, "StemcellAutomation.zip")
	outFile := filepath.Join(goPath, "src", buildPackage, "assets", "stemcell_automation.go")
	goBindataCommand := fmt.Sprintf("go-bindata -o %s -pkg assets -prefix %s %s", outFile, prefix, stemcellAutomation)
	goBindataCommandSlice := strings.Split(goBindataCommand, " ")
	goBinCommand := exec.Command(goBindataCommandSlice[0], goBindataCommandSlice[1:]...)
	goBinSession, err := gexec.Start(goBinCommand, stdout, stderr)
	if err != nil {
		return "", err
	}
	gomega.EventuallyWithOffset(1, goBinSession, 30*time.Second).Should(
		gexec.Exit(0),
		fmt.Sprintf("go-bindata command `%s` exited with exit code %d, stdout: %s, stderr: %s",
			goBindataCommand,
			goBinSession.ExitCode(),
			string(stdout.Bytes()),
			string(stderr.Bytes()),
		),
	)

	args := []string{
		"build",
		"-ldflags",
		fmt.Sprintf(`-X github.com/cloudfoundry-incubator/stembuild/version.Version=%s`, version),
		"-o",
		stembuildExecutable,
		buildPackage,
	}

	cmd := exec.Command("go", args...)

	session, err := gexec.Start(cmd, stdout, stderr)
	if err != nil {
		return "", err
	}
	gomega.EventuallyWithOffset(1, session, 30*time.Second).Should(
		gexec.Exit(0),
		fmt.Sprintf(
			"Build command was called with args: %v \n exited with exit code: %d, stdout: %s, stderr: %s",
			args,
			session.ExitCode(),
			string(stdout.Bytes()),
			string(stderr.Bytes()),
		),
	)

	return stembuildExecutable, err
}

func EnvMustExist(variableName string) string {
	result := os.Getenv(variableName)
	if result == "" {
		ginkgo.Fail(fmt.Sprintf("%s must be set", variableName))
	}

	return result
}
