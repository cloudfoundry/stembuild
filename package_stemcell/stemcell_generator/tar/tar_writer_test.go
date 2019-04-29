package tar_test

import (
	archiveTar "archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/stemcell_generator/tar"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TarWriter", func() {
	Describe("Write", func() {
		var (
			workingDir string
		)

		BeforeEach(func() {
			tmpDir := os.TempDir()
			var err error
			workingDir, err = ioutil.TempDir(tmpDir, "TarWriterTest")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			os.RemoveAll(workingDir)
		})

		It("should not fail", func() {
			w := tar.NewTarWriter()
			fakeFile := bytes.NewReader([]byte{})

			err := os.Chdir(workingDir)
			Expect(err).NotTo(HaveOccurred())

			err = w.Write("some-file", fakeFile)

			Expect(err).NotTo(HaveOccurred())
		})

		It("creates a file with the filename", func() {
			w := tar.NewTarWriter()
			fakeFile := bytes.NewReader([]byte{})

			err := os.Chdir(workingDir)
			Expect(err).NotTo(HaveOccurred())

			err = w.Write("some-file", fakeFile)

			Expect(err).NotTo(HaveOccurred())

			tarBall := filepath.Join(workingDir, "some-file")
			Expect(tarBall).To(BeAnExistingFile())
		})

		It("tars and zips the given readers", func() {
			w := tar.NewTarWriter()
			expectedContents := []string{"file1 content", "file2 content"}
			fakeFile1 := bytes.NewReader([]byte(expectedContents[0]))
			fakeFile2 := bytes.NewReader([]byte(expectedContents[1]))

			err := os.Chdir(workingDir)
			Expect(err).NotTo(HaveOccurred())

			err = w.Write("some-zipped-tar", fakeFile1, fakeFile2)

			var fileReader, _ = os.OpenFile("some-zipped-tar", os.O_RDONLY, 0777)

			gzr, err := gzip.NewReader(fileReader)
			Expect(err).ToNot(HaveOccurred())
			defer gzr.Close()
			tarfileReader := archiveTar.NewReader(gzr)
			var actualContents []string
			for {
				_, err := tarfileReader.Next()
				if err == io.EOF {
					break
				}

				Expect(err).NotTo(HaveOccurred())
				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(tarfileReader)
				if err != nil {
					break
				}
				actualContents = append(actualContents, buf.String())
			}
			Expect(len(actualContents)).To(Equal(2))
		})
	})
})
