package version_test

import (
	"github.com/cloudfoundry-incubator/stembuild/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetVersion", func() {
	It("should return a properly formatted version number", func() {
		os, stemcellVersion := version.GetVersions("1803.123.13")
		Expect(os).To(Equal("1803"))
		Expect(stemcellVersion).To(Equal("1803.123"))
	})

	It("should return 2016 as OS if given version is 1709", func() {
		os, stemcellVersion := version.GetVersions("1709.123.13")
		Expect(os).To(Equal("2016"))
		Expect(stemcellVersion).To(Equal("1709.123"))
	})

	It("should return 2012R2 as OS if given version is 1200", func() {
		os, stemcellVersion := version.GetVersions("1200.123.13")
		Expect(os).To(Equal("2012R2"))
		Expect(stemcellVersion).To(Equal("1200.123"))
	})
})
