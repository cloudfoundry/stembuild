package tar_test

import (
	archiveTar "archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/stembuild/package_stemcell/stemcell_generator/tar"
	"github.com/cloudfoundry/stembuild/package_stemcell/stemcell_generator/tar/tarfakes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TarWriter", func() {
	Describe("Write", func() {
		var (
			originalWorkingDir string
			workingDir         string
			fakeTarable        *tarfakes.FakeTarable
		)

		BeforeEach(func() {
			var err error
			originalWorkingDir, err = os.Getwd()
			Expect(err).NotTo(HaveOccurred())

			// Revert to manual cleanup which fails non-catastrophically on windows
			//workingDir = GinkgoT().TempDir() // automatically cleaned up
			workingDir, err = os.MkdirTemp(os.TempDir(), "TarWriterTest")
			Expect(err).NotTo(HaveOccurred())

			fakeFile := bytes.NewReader([]byte{})
			fakeTarable = &tarfakes.FakeTarable{}
			fakeTarable.ReadStub = fakeFile.Read
			fakeTarable.SizeStub = fakeFile.Size
			fakeTarable.NameReturns("some-file")
		})

		AfterEach(func() {
			Expect(os.Chdir(originalWorkingDir)).To(Succeed())

			// TODO: remove once GinkgoT().TempDir() is safe on windows
			err := os.RemoveAll(workingDir)
			if err != nil {
				By(fmt.Sprintf("removing '%s' failed: %s", workingDir, err))
			}
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
			fakeTarable1.NameReturns("first-file")

			fakeTarable2.ReadStub = fakeFile2.Read
			fakeTarable2.SizeStub = fakeFile2.Size
			fakeTarable2.NameReturns("second-file")

			err := os.Chdir(workingDir)
			Expect(err).NotTo(HaveOccurred())

			err = w.Write("some-zipped-tar", fakeTarable1, fakeTarable2) //nolint:ineffassign,staticcheck
			Expect(err).NotTo(HaveOccurred())

			var fileReader, _ = os.OpenFile("some-zipped-tar", os.O_RDONLY, 0777)

			gzr, err := gzip.NewReader(fileReader)
			Expect(err).ToNot(HaveOccurred())
			defer func() { _ = gzr.Close() }()

			tarReader := archiveTar.NewReader(gzr)
			var actualContents []string
			var actualFilenames []string
			for {
				header, err := tarReader.Next()
				if errors.Is(err, io.EOF) {
					break
				}
				Expect(err).NotTo(HaveOccurred())
				actualFilenames = append(actualFilenames, header.Name)
				Expect(header.Mode).To(Equal(int64(os.FileMode(0644))))

				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(tarReader)
				if err != nil {
					break
				}
				actualContents = append(actualContents, buf.String())
			}

			Expect(actualContents).To(ConsistOf(expectedContents))
			Expect(actualFilenames).To(ConsistOf("first-file", "second-file"))
		})
	})
})
