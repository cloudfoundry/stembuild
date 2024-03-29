package packagers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ bool = Describe("Packager Utility", func() {
	Context("TarGenerator", func() {
		var sourceDir string
		var destinationDir string

		BeforeEach(func() {
			sourceDir, _ = os.MkdirTemp(os.TempDir(), "packager-utility-test-source")
			destinationDir, _ = os.MkdirTemp(os.TempDir(), "packager-utility-test-destination")
		})

		AfterEach(func() {
			os.RemoveAll(sourceDir)      //nolint:errcheck
			os.RemoveAll(destinationDir) //nolint:errcheck
		})

		It("should tar all files inside provided folder and return its sha1", func() {
			err := os.WriteFile(filepath.Join(sourceDir, "file1"), []byte("file1 content\n"), 0777)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(filepath.Join(sourceDir, "file2"), []byte("file2 content\n"), 0777)
			Expect(err).NotTo(HaveOccurred())
			fileContentMap := make(map[string]string)
			fileContentMap["file1"] = "file1 content\n"
			fileContentMap["file2"] = "file2 content\n"

			tarball := filepath.Join(destinationDir, "tarball")

			sha1Sum, err := TarGenerator(tarball, sourceDir)

			Expect(err).NotTo(HaveOccurred())

			_, err = os.Stat(tarball)
			Expect(err).NotTo(HaveOccurred())
			var fileReader, _ = os.OpenFile(tarball, os.O_RDONLY, 0777)

			gzr, err := gzip.NewReader(fileReader)
			Expect(err).ToNot(HaveOccurred())
			defer gzr.Close()
			tarfileReader := tar.NewReader(gzr)
			count := 0
			for {
				header, err := tarfileReader.Next()
				if err == io.EOF {
					break
				}
				count++
				Expect(err).NotTo(HaveOccurred())
				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(tarfileReader)
				if err != nil {
					break
				}
				Expect(fileContentMap[header.Name]).To(Equal(buf.String()))
			}
			Expect(count).To(Equal(2))

			tarballFile, err := os.Open(tarball) //nolint:ineffassign,staticcheck
			defer tarballFile.Close()            //nolint:staticcheck
			expectedSha1 := sha1.New()
			io.Copy(expectedSha1, tarballFile) //nolint:errcheck

			sum := fmt.Sprintf("%x", expectedSha1.Sum(nil))
			Expect(sha1Sum).To(Equal(sum))
		})
	})

	Context("CreateManifest", func() {
		It("Creates a manifest correctly", func() {
			expectedManifest := `---
name: bosh-vsphere-esxi-windows1-go_agent
version: 'version'
api_version: 3
sha1: sha1sum
operating_system: windows1
cloud_properties:
  infrastructure: vsphere
  hypervisor: esxi
stemcell_formats:
- vsphere-ovf
- vsphere-ova
`
			result := CreateManifest("1", "version", "sha1sum")
			Expect(result).To(Equal(expectedManifest))
		})
	})

	Context("StemcellFileName", func() {
		It("formats a file name appropriately", func() {
			expectedName := "bosh-stemcell-1200.1-vsphere-esxi-windows2012R2-go_agent.tgz"
			Expect(StemcellFilename("1200.1", "2012R2")).To(Equal(expectedName))
		})
	})
})
