package commandparser_test

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"os"

	. "github.com/cloudfoundry-incubator/stembuild/commandparser"
	"github.com/cloudfoundry-incubator/stembuild/commandparser/commandparserfakes"

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

		It("stores the value of stemcell version", func() {
			err := f.Parse(longformArgs)
			Expect(err).ToNot(HaveOccurred())
			Expect(ConstrCmd.GetStemcellVersion()).To(Equal("1803.45"))
		})

		It("stores the value of a WinRM user", func() {
			err := f.Parse(longformArgs)
			Expect(err).ToNot(HaveOccurred())
			Expect(ConstrCmd.GetWinRMUser()).To(Equal("Admin"))
		})

		It("stores the value of a WinRM password ", func() {
			err := f.Parse(longformArgs)
			Expect(err).ToNot(HaveOccurred())
			Expect(ConstrCmd.GetWinRMPwd()).To(Equal("some_password"))
		})

		It("stores the value of the a WinRM IP", func() {
			err := f.Parse(longformArgs)
			Expect(err).ToNot(HaveOccurred())
			Expect(ConstrCmd.GetWinRMIp()).To(Equal("10.0.0.5"))
		})

	})

	Describe("Execute", func() {

		var f *flag.FlagSet
		var gf *GlobalFlags
		var ConstrCmd ConstructCmd
		var emptyContext context.Context

		var fakeFactory *commandparserfakes.FakeIVMConstructFactory
		var fakeVMPreparer *commandparserfakes.FakeVMPreparer

		BeforeEach(func() {
			f = flag.NewFlagSet("test", flag.ExitOnError)
			gf = &GlobalFlags{true, true, true}

			fakeFactory = &commandparserfakes.FakeIVMConstructFactory{}
			fakeVMPreparer = &commandparserfakes.FakeVMPreparer{}
			fakeFactory.GetVMPreparerReturns(fakeVMPreparer)

			ConstrCmd.SetFlags(f)
			ConstrCmd = NewConstructCmd(fakeFactory)
			ConstrCmd.GlobalFlags = gf
			emptyContext = context.Background()
			os.Create("LGPO.zip")
		})

		AfterSuite(func() {
			os.Remove("LGPO.zip")
		})

		It("should execute the construct VM command", func() {
			args := []string{"-stemcell-version", "1803.45",
				"-winrm-ip", "10.0.0.5",
				"-winrm-username", "Admin",
				"-winrm-password", "some_password",
			}
			err := f.Parse(args)
			Expect(err).ToNot(HaveOccurred())

			exitStatus := ConstrCmd.Execute(emptyContext, f)

			Expect(exitStatus).To(Equal(subcommands.ExitSuccess))
		})
	})

})
