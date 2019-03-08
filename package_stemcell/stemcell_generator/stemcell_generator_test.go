package stemcell_generator_test

import (
	"bytes"
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/stemcell_generator"
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/stemcell_generator/generatorfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("StemcellGenerator", func() {
	Describe("Generate", func() {
		var (
			stemcellGenerator *stemcell_generator.StemcellGenerator
			manifestGenerator *generatorfakes.FakeManifestGenerator
			fakeImage         io.Reader
		)

		BeforeEach(func() {
			fakeImage = bytes.NewReader([]byte{})
			manifestGenerator = &generatorfakes.FakeManifestGenerator{}
			stemcellGenerator = stemcell_generator.NewStemcellGenerator(manifestGenerator)
		})

		It("generates a manifest", func() {
			err := stemcellGenerator.Generate(fakeImage)

			Expect(err).NotTo(HaveOccurred())
			Expect(manifestGenerator.GenerateCallCount()).To(Equal(1))

			args := manifestGenerator.GenerateArgsForCall(0)
			Expect(args).To(Equal(fakeImage))
		})


	})
})
