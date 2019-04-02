package construct_test

import (
	"github.com/cloudfoundry-incubator/stembuild/construct"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Messenger", func() {
	var buf *gbytes.Buffer

	BeforeEach(func() {
		buf = gbytes.NewBuffer()
	})

	Describe("EnableWinRMStarted", func() {
		It("writes the message to the writer", func() {

			m := construct.NewMessenger(buf)
			m.EnableWinRMStarted()

			Expect(buf).To(gbytes.Say("\nAttempting to enable WinRM on the guest vm..."))
		})
	})

	Describe("EnableWinRMSucceeded", func() {
		It("writes the message to the writer", func() {

			m := construct.NewMessenger(buf)
			m.EnableWinRMSucceeded()

			Expect(buf).To(gbytes.Say("WinRm enabled on the guest VM\n"))
		})
	})

	Describe("WinRM", func() {
		It("writes both WinRM messages on one line", func() {

			m := construct.NewMessenger(buf)
			m.EnableWinRMStarted()
			m.EnableWinRMSucceeded()

			Expect(buf).To(gbytes.Say("Attempting to enable WinRM on the guest vm...WinRm enabled on the guest VM"))
		})
	})

	Describe("ValidateVMConnectionStarted", func() {
		It("writes the message to the writer", func() {

			m := construct.NewMessenger(buf)
			m.ValidateVMConnectionStarted()

			Expect(buf).To(gbytes.Say("\nValidating connection to vm..."))
		})
	})

	Describe("ValidateVMConnectionSucceeded", func() {
		It("writes the message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.ValidateVMConnectionSucceeded()

			Expect(buf).To(gbytes.Say("succeeded.\n"))
		})
	})

	Describe("ValidateVMConnection", func() {
		It("writes both validate vm connection messages on one line", func() {
			m := construct.NewMessenger(buf)
			m.ValidateVMConnectionStarted()
			m.ValidateVMConnectionSucceeded()

			Expect(buf).To(gbytes.Say("Validating connection to vm...succeeded."))
		})
	})

	Describe("CreateProvisionDirStarted", func() {
		It("writes the message to the writer", func() {

			m := construct.NewMessenger(buf)
			m.CreateProvisionDirStarted()

			Expect(buf).To(gbytes.Say("\nCreating provision dir on target VM..."))
		})
	})

	Describe("CreateProvisionDirSucceeded", func() {
		It("writes the message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.CreateProvisionDirSucceeded()

			Expect(buf).To(gbytes.Say("succeeded.\n"))
		})
	})

	Describe("CreateProvisionDir", func() {
		It("writes both messages on one line", func() {
			m := construct.NewMessenger(buf)
			m.CreateProvisionDirStarted()
			m.CreateProvisionDirSucceeded()

			Expect(buf).To(gbytes.Say("\nCreating provision dir on target VM...succeeded.\n"))
		})
	})
})
