package stemcell_generator_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/cloudfoundry/stembuild/package_stemcell/stemcell_generator"
	fakes "github.com/cloudfoundry/stembuild/package_stemcell/stemcell_generator/stemcell_generatorfakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StemcellGenerator", func() {
	Describe("Generate", func() {
		var (
			stemcellGenerator *stemcell_generator.StemcellGenerator
			manifestGenerator *fakes.FakeManifestGenerator
			fileNameGenerator *fakes.FakeFileNameGenerator
			tarWriter         *fakes.FakeTarWriter
			fakeImage         io.Reader
		)

		BeforeEach(func() {
			fakeImage = bytes.NewReader([]byte{})
			manifestGenerator = &fakes.FakeManifestGenerator{}
			fileNameGenerator = &fakes.FakeFileNameGenerator{}
			tarWriter = &fakes.FakeTarWriter{}
			stemcellGenerator = stemcell_generator.NewStemcellGenerator(manifestGenerator, fileNameGenerator, tarWriter)
		})

		It("generates a manifest", func() {
			err := stemcellGenerator.Generate(fakeImage)

			Expect(err).NotTo(HaveOccurred())
			Expect(manifestGenerator.ManifestCallCount()).To(Equal(1))

			args := manifestGenerator.ManifestArgsForCall(0)
			Expect(args).To(Equal(fakeImage))
		})

		It("returns an error when manifest generation fails", func() {
			manifestGenerator.ManifestReturns(nil, errors.New("some manifest error"))

			err := stemcellGenerator.Generate(fakeImage)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("failed to generate stemcell manifest: some manifest error"))
		})

		It("generates a filename", func() {
			err := stemcellGenerator.Generate(fakeImage)

			Expect(err).NotTo(HaveOccurred())
			Expect(fileNameGenerator.FileNameCallCount()).To(Equal(1))
		})

		It("should generate a tarball", func() {
			expectedFileName := "the-file.tgz"
			expectedManifest := bytes.NewReader([]byte("manifest"))

			fileNameGenerator.FileNameReturns(expectedFileName)
			manifestGenerator.ManifestReturns(expectedManifest, nil)

			stemcellGenerator.Generate(fakeImage) //nolint:errcheck
			Expect(tarWriter.WriteCallCount()).To(Equal(1))

			actualFileName, objects := tarWriter.WriteArgsForCall(0)

			Expect(actualFileName).To(Equal(expectedFileName))

			Expect(objects).To(ConsistOf(expectedManifest, fakeImage))
		})

		It("should return an error when tar writer fails", func() {
			tarWriterError := errors.New("some tar writer error")
			tarWriter.WriteReturns(tarWriterError)

			err := stemcellGenerator.Generate(fakeImage)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fmt.Sprintf("failed to generate stemcell tarball: %s", tarWriterError)))
		})
	})
})
