package packagers

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func WriteManifest(manifestContents, manifestPath string) error {

	manifestPath = filepath.Join(manifestPath, "stemcell.MF")

	f, err := os.OpenFile(manifestPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("creating stemcell.MF (%s): %s", manifestPath, err)
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, manifestContents); err != nil {
		os.Remove(manifestPath)
		return fmt.Errorf("writing stemcell.MF (%s): %s", manifestPath, err)
	}
	return nil
}

func CreateManifest(osVersion, version, sha1sum string) string {
	const format = `---
name: bosh-vsphere-esxi-windows%[1]s-go_agent
version: '%[2]s'
sha1: %[3]s
operating_system: windows%[1]s
cloud_properties:
  infrastructure: vsphere
  hypervisor: esxi
stemcell_formats:
- vsphere-ovf
- vsphere-ova
`
	return fmt.Sprintf(format, osVersion, version, sha1sum)

}

func TarGenerator(destinationfileName string, sourceDirName string) (string, error) {

	sourcedir, err := os.Open(sourceDirName)
	if err != nil {
		return "", errors.New(fmt.Sprintf("unable to open %s", sourceDirName))
	}
	defer sourcedir.Close()

	// get list of files
	files, err := sourcedir.Readdir(0)
	if err != nil {
		return "", errors.New(fmt.Sprintf("unable to list files in %s", sourceDirName))
	}

	// create tar file
	destinationFile, err := os.Create(destinationfileName)
	if err != nil {
		return "", errors.New(fmt.Sprintf("unable to create destination file with name %s", destinationfileName))
	}
	defer destinationFile.Close()

	sha1Hash := sha1.New()
	gzw := gzip.NewWriter(io.MultiWriter(destinationFile, sha1Hash))
	tarfileWriter := tar.NewWriter(gzw)

	for _, fileInfo := range files {

		if fileInfo.IsDir() {
			continue
		}

		file, err := os.Open(sourcedir.Name() + string(filepath.Separator) + fileInfo.Name())
		if err != nil {
			return "", errors.New(fmt.Sprintf("unable to open files in %s", sourceDirName))
		}
		defer file.Close()

		// prepare the tar header
		header := new(tar.Header)
		header.Name = fileInfo.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()

		err = tarfileWriter.WriteHeader(header)
		if err != nil {
			return "", errors.New("unable to write to header of destination tar file")
		}

		_, err = io.Copy(tarfileWriter, file)
		if err != nil {
			return "", errors.New("unable to write contents to destination tar file")
		}
	}

	//Shouldn't be a deferred call as closing the tar writer flushes padding and writes footer which impacts the sha1sum

	err = tarfileWriter.Close()
	if err != nil {
		return "", errors.New("unable to close tar file")
	}
	err = gzw.Close()
	if err != nil {
		return "", errors.New("unable to close tar file (gzip)")
	}

	return fmt.Sprintf("%x", sha1Hash.Sum(nil)), nil
}

func StemcellFilename(version, os string) string {
	return fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz",
		version, os)
}
