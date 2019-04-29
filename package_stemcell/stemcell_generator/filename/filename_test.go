package filename_test

import (
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/stemcell_generator/filename"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filename", func() {

	It("should return the expected filename", func() {
		filenameGenerator := filename.NewFilenameGenerator("2016", "1709.999")

		f := filenameGenerator.Filename()

		Expect(f).To(Equal("bosh-stemcell-1709.999-vsphere-esxi-windows2016-go_agent.tgz"))
	})
})
