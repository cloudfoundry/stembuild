package stemcell_generator_test

import (
	"bytes"
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/stemcell_generator"
	fakes "github.com/cloudfoundry-incubator/stembuild/package_stemcell/stemcell_generator/stemcell_generatorfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("StemcellGenerator", func() {
	Describe("Generate", func() {
		var (
			stemcellGenerator *stemcell_generator.StemcellGenerator
			manifestGenerator *fakes.FakeManifestGenerator
			fileNameGenerator *fakes.FakeFileNameGenerator
			tarWriter *fakes.FakeTarWriter
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

		It("generates a filename", func() {
			err := stemcellGenerator.Generate(fakeImage)

			Expect(err).NotTo(HaveOccurred())
			Expect(fileNameGenerator.FileNameCallCount()).To(Equal(1))
		})
		It("should generate a tarball", func(){
			expectedFileName := "the-file.tgz"
			expectedManifest := bytes.NewReader([]byte("manifest"))

			fileNameGenerator.FileNameReturns(expectedFileName)
			manifestGenerator.ManifestReturns(expectedManifest, nil)

			stemcellGenerator.Generate(fakeImage)
			Expect(tarWriter.WriteCallCount()).To(Equal(1))

			actualFileName, objects := tarWriter.WriteArgsForCall(0)

			Expect(actualFileName).To(Equal(expectedFileName))

			Expect(objects).To(ConsistOf(expectedManifest, fakeImage))
		})

	})
})
