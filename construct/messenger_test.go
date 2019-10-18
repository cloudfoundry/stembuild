package construct_test

import (
	"fmt"

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

	Describe("Enable WinRM messages", func() {
		It("writes the started message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.EnableWinRMStarted()

			Expect(buf).To(gbytes.Say("\nAttempting to enable WinRM on the guest vm..."))
		})

		It("writes the succeeded message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.EnableWinRMSucceeded()

			Expect(buf).To(gbytes.Say("WinRm enabled on the guest VM\n"))
		})

		It("writes both WinRM messages on one line", func() {
			m := construct.NewMessenger(buf)
			m.EnableWinRMStarted()
			m.EnableWinRMSucceeded()

			Expect(buf).To(gbytes.Say("Attempting to enable WinRM on the guest vm...WinRm enabled on the guest VM"))
		})
	})

	Describe("Validate VM connection messages", func() {
		It("writes the started message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.ValidateVMConnectionStarted()

			Expect(buf).To(gbytes.Say("\nValidating connection to vm..."))
		})

		It("writes the succeeded message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.ValidateVMConnectionSucceeded()

			Expect(buf).To(gbytes.Say("succeeded.\n"))
		})

		It("writes both validate vm connection messages on one line", func() {
			m := construct.NewMessenger(buf)
			m.ValidateVMConnectionStarted()
			m.ValidateVMConnectionSucceeded()

			Expect(buf).To(gbytes.Say("Validating connection to vm...succeeded."))
		})
	})

	Describe("Create provision directory messages", func() {
		It("writes the started message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.CreateProvisionDirStarted()

			Expect(buf).To(gbytes.Say("\nCreating provision dir on target VM..."))
		})

		It("writes the succeeded message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.CreateProvisionDirSucceeded()

			Expect(buf).To(gbytes.Say("succeeded.\n"))
		})

		It("writes both messages on one line", func() {
			m := construct.NewMessenger(buf)
			m.CreateProvisionDirStarted()
			m.CreateProvisionDirSucceeded()

			Expect(buf).To(gbytes.Say("\nCreating provision dir on target VM...succeeded.\n"))
		})
	})

	Describe("Upload artifacts messages", func() {
		It("writes the started message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.UploadArtifactsStarted()

			Expect(buf).To(gbytes.Say("\nTransferring ~20 MB to the Windows VM. Depending on your connection, the transfer may take 15-45 minutes\n"))
		})

		It("writes the succeeded message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.UploadArtifactsSucceeded()

			Expect(buf).To(gbytes.Say("\nAll files have been uploaded.\n"))
		})
	})

	Describe("Extract artifact messages", func() {
		It("writes the started message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.ExtractArtifactsStarted()

			Expect(buf).To(gbytes.Say("\nExtracting artifacts..."))
		})

		It("writes the succeeded message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.ExtractArtifactsSucceeded()

			Expect(buf).To(gbytes.Say("succeeded.\n"))
		})

		It("writes both messages on one line", func() {
			m := construct.NewMessenger(buf)
			m.ExtractArtifactsStarted()
			m.ExtractArtifactsSucceeded()

			Expect(buf).To(gbytes.Say("\nExtracting artifacts...succeeded.\n"))
		})
	})

	Describe("Execute setup script messages", func() {
		It("writes the started message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.ExecuteScriptStarted()

			Expect(buf).To(gbytes.Say("\nExecuting setup script...\n"))
		})

		It("writes the succeeded message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.ExecuteScriptSucceeded()

			Expect(buf).To(gbytes.Say("\nFinished executing setup script.\n"))
		})

		It("writes the restarting message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.WinRMDisconnectedForReboot()

			Expect(buf).To(gbytes.Say("WinRM has been disconnected so the VM can reboot. Preparing the VM to be shutdown."))

		})
	})

	Describe("Upload file messages", func() {
		It("writes the started message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.UploadFileStarted("some artifact")

			Expect(buf).To(gbytes.Say("\tUploading some artifact to target VM..."))
		})

		It("writes the succeeded message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.UploadFileSucceeded()

			Expect(buf).To(gbytes.Say("succeeded.\n"))
		})

		It("writes both messages on one line", func() {
			m := construct.NewMessenger(buf)
			m.UploadFileStarted("some third artifact")
			m.UploadFileSucceeded()

			Expect(buf).To(gbytes.Say("Uploading some third artifact to target VM...succeed."))
		})
	})

	Describe("validate OS", func() {
		var matchingVersionWarning = "Ensure the version of the stemcell you're trying to build matches the corresponding base ISO you're using.\n" +
			"For example: If you're building 2019.x, then you should be using 'Windows Server 2019' only"
		It("writes the OS version file creation failed message to the writer", func() {
			errorMessage := "some error message"
			m := construct.NewMessenger(buf)
			m.OSVersionFileCreationFailed(errorMessage)
			Expect(buf).To(gbytes.Say(fmt.Sprintf("Warning: OS Version file creation failed:\n%s\n%s", matchingVersionWarning, errorMessage)))
		})

		It("writes the exit code retrieval failed message to the writer", func() {
			errorMessage := "some error message"
			m := construct.NewMessenger(buf)
			m.ExitCodeRetrievalFailed(errorMessage)
			Expect(buf).To(gbytes.Say(fmt.Sprintf("Warning: Failed to retrieve exit code for process to create OS Version file:\n%s\n%s", matchingVersionWarning, errorMessage)))
		})

		It("writes the download file failed message to the writer", func() {
			errorMessage := "some error message"
			m := construct.NewMessenger(buf)
			m.DownloadFileFailed(errorMessage)
			Expect(buf).To(gbytes.Say(fmt.Sprintf("Warning: Failed to download OS Version file:\n%s\n%s", matchingVersionWarning, errorMessage)))
		})
	})

	Describe("Power state messages", func() {
		It("writes still running message with timestamp", func() {
			m := construct.NewMessenger(buf)
			m.RestartInProgress()

			dateRegexString := "[0-9-]*"
			timeRegexString := "[0-9:.]*"
			messageString := "still preparing VM...\n"
			logLineRegex := fmt.Sprintf("\\[.*%s\\s*%s.*\\]\\s*%s", dateRegexString, timeRegexString, messageString)

			Expect(buf).To(gbytes.Say(logLineRegex))
		})

		It("writes the shutdown message to the writer", func() {
			m := construct.NewMessenger(buf)
			m.ShutdownCompleted()

			Expect(buf).To(gbytes.Say("VM has now been shutdown. Run `stembuild package` to finish building the stemcell.\n"))
		})

	})
})
