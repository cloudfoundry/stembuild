package helpers

import (
	"bytes"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pivotal-cf-experimental/stembuild/stembuildoptions"
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

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func StringFromManifest(fileTemplate string, manifestStruct stembuildoptions.StembuildOptions) (string, error) {
	t, err := template.New("manifest template").Parse(fileTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, manifestStruct); err != nil {
		return "", err
	}
	return buf.String(), nil
}

const ManifestTemplate = `---
version: "{{.Version}}"
vhd_file: "{{.VHDFile}}"
patch_file: "{{.PatchFile}}"
os_version: "{{.OSVersion}}"
output_dir: "{{.OutputDir}}"
`
