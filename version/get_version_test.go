package version_test

import (
	"github.com/cloudfoundry-incubator/stembuild/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version Utilities", func() {
	Describe("GetVersions", func() {
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

	Describe("GetOSVersionFromBuildNumber", func() {
		It("returns 2019 given build number 17763", func() {
			Expect(version.GetOSVersionFromBuildNumber("17763")).To(Equal("2019"))
		})
		It("returns 1803 given build number 17134", func() {
			Expect(version.GetOSVersionFromBuildNumber("17134")).To(Equal("1803"))
		})
		It("returns empty string if given build number is wrong", func() {
			Expect(version.GetOSVersionFromBuildNumber("random")).To(Equal(""))
		})
	})
})
