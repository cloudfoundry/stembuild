package packagers_test

import (
	"bytes"
	"errors"

	"github.com/cloudfoundry/stembuild/package_stemcell/packagers"
	"github.com/cloudfoundry/stembuild/package_stemcell/packagers/packagersfakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Packager", func() {
	Describe("Package", func() {

		var (
			source            *packagersfakes.FakeSource
			stemcellGenerator *packagersfakes.FakeStemcellGenerator
			packager          *packagers.Packager
		)

		BeforeEach(func() {
			source = &packagersfakes.FakeSource{}
			stemcellGenerator = &packagersfakes.FakeStemcellGenerator{}
			packager = packagers.NewPackager(source, stemcellGenerator)
		})

		It("returns an error if ArtifactReader does", func() {

			source.ArtifactReaderReturns(nil, errors.New("bad thing"))
			err := packager.Package()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("packager failed to retrieve artifact: bad thing"))
		})

		It("returns an error if Generate does", func() {

			stemcellGenerator.GenerateReturns(errors.New("other bad thing"))

			err := packager.Package()

			Expect(err).To(MatchError("packager failed to generate stemcell: other bad thing"))
		})

		It("uses source object to generate stemcell", func() {
			fakeIoReader := bytes.NewReader([]byte{})
			source.ArtifactReaderReturns(fakeIoReader, nil)

			err := packager.Package()

			Expect(err).NotTo(HaveOccurred())

			Expect(source.ArtifactReaderCallCount()).To(Equal(1))
			Expect(stemcellGenerator.GenerateCallCount()).To(Equal(1))

			argsForFirstCall := stemcellGenerator.GenerateArgsForCall(0)
			Expect(argsForFirstCall).To(BeIdenticalTo(fakeIoReader))
		})
	})
})
