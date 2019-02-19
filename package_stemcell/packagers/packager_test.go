package packagers_test

import (
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/packagers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)



var _ = FDescribe("Packager", func() {
	Describe("Package", func() {
		It("doesn't return an error", func() {
			packager := &packagers.Packager{}

			err := packager.Package()
			Expect(err).NotTo(HaveOccurred())
		})
		It("calls ArtifactReader method on source object", func(){
			source := FakeSource
		})
	})
})
