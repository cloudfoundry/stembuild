package tar_test

import (
	archiveTar "archive/tar"
	"bytes"
	"compress/gzip"
	"github.com/cloudfoundry/stembuild/package_stemcell/stemcell_generator/tar/tarfakes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/stembuild/package_stemcell/stemcell_generator/tar"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TarWriter", func() {
	Describe("Write", func() {
		var (
			workingDir  string
			fakeTarable *tarfakes.FakeTarable
		)

		BeforeEach(func() {
			tmpDir := os.TempDir()
			var err error
			workingDir, err = ioutil.TempDir(tmpDir, "TarWriterTest")
			Expect(err).NotTo(HaveOccurred())

			fakeFile := bytes.NewReader([]byte{})
			fakeTarable = &tarfakes.FakeTarable{}
			fakeTarable.ReadStub = fakeFile.Read
			fakeTarable.SizeStub = fakeFile.Size
			fakeTarable.NameReturns("some-file")
		})

		AfterEach(func() {
			os.RemoveAll(workingDir)
		})

		It("should not fail", func() {
			w := tar.NewTarWriter()

			err := os.Chdir(workingDir)
			Expect(err).NotTo(HaveOccurred())

			err = w.Write("some-file", fakeTarable)

			Expect(err).NotTo(HaveOccurred())
		})

		It("creates a file with the filename", func() {
			w := tar.NewTarWriter()

			err := os.Chdir(workingDir)
			Expect(err).NotTo(HaveOccurred())

			err = w.Write("some-file", fakeTarable)

			Expect(err).NotTo(HaveOccurred())

			tarBall := filepath.Join(workingDir, "some-file")
			Expect(tarBall).To(BeAnExistingFile())
		})

		It("tars and zips the given readers", func() {
			w := tar.NewTarWriter()
			expectedContents := []string{"file1 content", "file2 slightly longer content"}
			fakeFile1 := bytes.NewReader([]byte(expectedContents[0]))
			fakeFile2 := bytes.NewReader([]byte(expectedContents[1]))
			fakeTarable1 := &tarfakes.FakeTarable{}
			fakeTarable2 := &tarfakes.FakeTarable{}

			fakeTarable1.ReadStub = fakeFile1.Read
			fakeTarable1.SizeStub = fakeFile1.Size
			fakeTarable1.NameReturns("firstfile")

			fakeTarable2.ReadStub = fakeFile2.Read
			fakeTarable2.SizeStub = fakeFile2.Size
			fakeTarable2.NameReturns("secondfile")

			err := os.Chdir(workingDir)
			Expect(err).NotTo(HaveOccurred())

			err = w.Write("some-zipped-tar", fakeTarable1, fakeTarable2)

			var fileReader, _ = os.OpenFile("some-zipped-tar", os.O_RDONLY, 0777)

			gzr, err := gzip.NewReader(fileReader)
			Expect(err).ToNot(HaveOccurred())
			defer gzr.Close()
			tarfileReader := archiveTar.NewReader(gzr)
			var actualContents []string
			var actualFilenames []string
			for {
				header, err := tarfileReader.Next()
				if err == io.EOF {
					break
				}
				Expect(err).NotTo(HaveOccurred())
				actualFilenames = append(actualFilenames, header.Name)
				Expect(header.Mode).To(Equal(int64(os.FileMode(0644))))

				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(tarfileReader)
				if err != nil {
					break
				}
				actualContents = append(actualContents, buf.String())
			}

			Expect(actualContents).To(ConsistOf(expectedContents))
			Expect(actualFilenames).To(ConsistOf("firstfile", "secondfile"))

		})
	})
})
