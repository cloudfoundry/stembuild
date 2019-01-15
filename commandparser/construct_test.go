package commandparser_test

import (
	"flag"
	. "github.com/cloudfoundry-incubator/stembuild/commandparser"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("construct", func() {
	// Focus of this test is not to test the Flags.Parse functionality as much
	// as to test that the command line flags values are stored in the expected
	// struct variables. This adds a bit of protection when renaming flag parameters.
	Describe("SetFlags", func() {

		var f *flag.FlagSet
		var ConstrCmd *ConstructCmd

		BeforeEach(func() {
			f = flag.NewFlagSet("test", flag.ExitOnError)
			ConstrCmd = &ConstructCmd{}
			ConstrCmd.SetFlags(f)
		})

		var longformArgs = []string{"-stemcell-version", "1803.45",
			"-winrm-ip", "10.0.0.5",
			"-winrm-username", "Admin",
			"-winrm-password", "some_password",
		}
		var shortformArgs = []string{"-s", "1803.45",
			"-ip", "10.0.0.5",
			"-u", "Admin",
			"-p", "some_password",
		}

		Context("a stemcell version is specified as a flag parameter", func() {
			It("stores the value when using the long form", func() {
				err := f.Parse(longformArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(ConstrCmd.GetStemcellVersion()).To(Equal("1803.45"))
			})

			It("stores the value when using the short form", func() {
				err := f.Parse(shortformArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(ConstrCmd.GetStemcellVersion()).To(Equal("1803.45"))
			})

		})

		Context("a WinRM user is specified as a flag parameter", func() {
			It("stores the value when using the long form", func() {
				err := f.Parse(longformArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(ConstrCmd.GetWinRMUser()).To(Equal("Admin"))
			})

			It("stores the value when using the short form", func() {
				err := f.Parse(shortformArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(ConstrCmd.GetWinRMUser()).To(Equal("Admin"))
			})

		})

		Context("a WinRM password is specified as a flag parameter", func() {
			It("stores the value when using the long form", func() {
				err := f.Parse(longformArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(ConstrCmd.GetWinRMPwd()).To(Equal("some_password"))
			})

			It("stores the value when using the short form", func() {
				err := f.Parse(shortformArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(ConstrCmd.GetWinRMPwd()).To(Equal("some_password"))
			})

		})

		Context("a WinRM IP is specified as a flag parameter", func() {
			It("stores the value when using the long form", func() {
				err := f.Parse(longformArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(ConstrCmd.GetWinRMIp()).To(Equal("10.0.0.5"))
			})

			It("stores the value when using the short form", func() {
				err := f.Parse(shortformArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(ConstrCmd.GetWinRMIp()).To(Equal("10.0.0.5"))
			})

		})

	})
})
