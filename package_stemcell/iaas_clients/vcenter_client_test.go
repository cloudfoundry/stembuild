package iaas_clients

import (
	"fmt"

	"github.com/cloudfoundry-incubator/stembuild/iaas_cli/iaas_clifakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VcenterClient", func() {
	var (
		runner                  *iaas_clifakes.FakeCliRunner
		username, password, url string
		vcenterClient           *VcenterClient
		credentialUrl           string
	)

	BeforeEach(func() {
		runner = &iaas_clifakes.FakeCliRunner{}
		username, password, url = "username", "password", "url"
		vcenterClient = NewVcenterClient(username, password, url, runner)
		credentialUrl = fmt.Sprintf("%s:%s@%s", username, password, url)
	})

	Context("ValidateCredentials", func() {
		It("When the login credentials are correct, login is successful", func() {
			expectedArgs := []string{"about", "-u", credentialUrl}

			runner.RunReturns(0)
			err := vcenterClient.ValidateCredentials()
			argsForRun := runner.RunArgsForCall(0)

			Expect(err).To(Not(HaveOccurred()))
			Expect(runner.RunCallCount()).To(Equal(1))
			Expect(argsForRun).To(Equal(expectedArgs))
		})

		It("When the login credentials are incorrect, login is a failure", func() {
			expectedArgs := []string{"about", "-u", credentialUrl}

			runner.RunReturns(1)
			err := vcenterClient.ValidateCredentials()
			argsForRun := runner.RunArgsForCall(0)

			Expect(err).To(HaveOccurred())
			Expect(runner.RunCallCount()).To(Equal(1))
			Expect(argsForRun).To(Equal(expectedArgs))
			Expect(err.Error()).To(Equal("invalid credentials"))
		})
	})

	Context("validateUrl", func() {
		It("When the url is valid, there is no error", func() {
			expectedArgs := []string{"about", "-u", url}

			runner.RunReturns(0)
			err := vcenterClient.ValidateUrl()
			argsForRun := runner.RunArgsForCall(0)

			Expect(err).To(Not(HaveOccurred()))
			Expect(runner.RunCallCount()).To(Equal(1))
			Expect(argsForRun).To(Equal(expectedArgs))
		})

		It("When the url is invalid, there is an error", func() {
			expectedArgs := []string{"about", "-u", url}

			runner.RunReturns(1)
			err := vcenterClient.ValidateUrl()
			argsForRun := runner.RunArgsForCall(0)

			Expect(err).To(HaveOccurred())
			Expect(runner.RunCallCount()).To(Equal(1))
			Expect(argsForRun).To(Equal(expectedArgs))
			Expect(err.Error()).To(Equal("invalid url"))
		})
	})

	Context("FindVM", func() {
		It("If the VM path is valid, and the VM is found", func() {
			expectedArgs := []string{"find", "-maxdepth=0", "-u", credentialUrl, "validVMPath"}
			runner.RunReturns(0)
			err := vcenterClient.FindVM("validVMPath")
			argsForRun := runner.RunArgsForCall(0)

			Expect(err).To(Not(HaveOccurred()))
			Expect(runner.RunCallCount()).To(Equal(1))
			Expect(argsForRun).To(Equal(expectedArgs))
		})

		It("If the VM path is valid, and the VM is found", func() {
			expectedArgs := []string{"find", "-maxdepth=0", "-u", credentialUrl, "validVMPath"}
			runner.RunReturns(1)
			err := vcenterClient.FindVM("validVMPath")
			argsForRun := runner.RunArgsForCall(0)

			Expect(err).To(HaveOccurred())
			Expect(runner.RunCallCount()).To(Equal(1))
			Expect(argsForRun).To(Equal(expectedArgs))
			Expect(err.Error()).To(Equal("invalid VM path"))
		})
	})

	Context("PrepareVM", func() {
		It("removes the ethernet device", func() {
			expectedArgs := []string{"device.remove", "-vm", "validVMPath", "ethernet-0", "-u", credentialUrl}
			runner.RunReturns(0)
			err := vcenterClient.PrepareVM("validVMPath")

			Expect(err).To(Not(HaveOccurred()))
			Expect(runner.RunCallCount()).To(Equal(2))

			argsForRun := runner.RunArgsForCall(0)
			Expect(argsForRun).To(Equal(expectedArgs))

		})

		It("returns an error if it was not able to remove ethernet-0 for some reason", func() {
			expectedArgs := []string{"device.remove", "-vm", "validVMPath", "ethernet-0", "-u", credentialUrl}
			runner.RunReturns(1)
			err := vcenterClient.PrepareVM("validVMPath")

			Expect(err).To(HaveOccurred())
			Expect(runner.RunCallCount()).To(Equal(1))

			argsForRun := runner.RunArgsForCall(0)
			Expect(argsForRun).To(Equal(expectedArgs))
			Expect(err.Error()).To(Equal("ethernet-0 could not be removed/not found"))
		})

		It("removes the virtual floppy device", func() {
			expectedArgs := []string{"device.remove", "-vm", "validVMPath", "floppy-8000", "-u", credentialUrl}
			runner.RunReturns(0)
			runner.RunReturnsOnCall(0, 0)
			runner.RunReturnsOnCall(1, 0)
			err := vcenterClient.PrepareVM("validVMPath")

			Expect(err).To(Not(HaveOccurred()))
			Expect(runner.RunCallCount()).To(Equal(2))

			argsForRun := runner.RunArgsForCall(1)
			Expect(argsForRun).To(Equal(expectedArgs))

		})

		It("returns an error if it was not able to remove ethernet-0 for some reason", func() {
			expectedArgs := []string{"device.remove", "-vm", "validVMPath", "floppy-8000", "-u", credentialUrl}
			runner.RunReturnsOnCall(0, 0)
			runner.RunReturnsOnCall(1, 1)
			err := vcenterClient.PrepareVM("validVMPath")

			Expect(err).To(HaveOccurred())
			Expect(runner.RunCallCount()).To(Equal(2))

			argsForRun := runner.RunArgsForCall(1)
			Expect(argsForRun).To(Equal(expectedArgs))
			Expect(err.Error()).To(Equal("floppy-8000 could not be removed/not found"))
		})
	})
})
