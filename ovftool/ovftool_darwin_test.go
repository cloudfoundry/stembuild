// +build darwin

package ovftool_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"

	"github.com/pivotal-cf-experimental/stembuild/ovftool"
)

var _ = Describe("ovftool darwin", func() {

	Context("homeDirectory", func() {
		var home = ""

		BeforeEach(func() {
			home = os.Getenv("HOME")
		})

		AfterEach(func() {
			os.Setenv("HOME", home)
		})

		It("returns value of HOME environment variable is set", func() {
			Expect(home).NotTo(Equal(""))

			testHome := ovftool.HomeDirectory()
			Expect(testHome).To(Equal(home))
		})

		It("returns user HOME directory if HOME environment variable is not set", func() {
			os.Unsetenv("HOME")

			testHome := ovftool.HomeDirectory()
			Expect(testHome).To(Equal(home))
		})

	})

	// OVFTool not unit tested as it is just a wrapper for findExecutable
	// with a array of root directories

})
