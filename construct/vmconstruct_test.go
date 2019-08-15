package construct_test

import (
	"context"
	"errors"

	"github.com/onsi/gomega/gbytes"

	. "github.com/cloudfoundry-incubator/stembuild/construct"
	"github.com/cloudfoundry-incubator/stembuild/construct/constructfakes"
	"github.com/cloudfoundry-incubator/stembuild/remotemanager/remotemanagerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("construct_helpers", func() {
	var (
		fakeRemoteManager *remotemanagerfakes.FakeRemoteManager
		vmConstruct       *VMConstruct
		fakeVcenterClient *constructfakes.FakeIaasClient
		fakeGuestManager  *constructfakes.FakeGuestManager
		fakeWinRMEnabler  *constructfakes.FakeWinRMEnabler
		fakeOSValidator   *constructfakes.FakeOSValidator
		fakeMessenger     *constructfakes.FakeConstructMessenger
	)

	BeforeEach(func() {
		fakeRemoteManager = &remotemanagerfakes.FakeRemoteManager{}
		fakeVcenterClient = &constructfakes.FakeIaasClient{}
		fakeGuestManager = &constructfakes.FakeGuestManager{}
		fakeWinRMEnabler = &constructfakes.FakeWinRMEnabler{}
		fakeOSValidator = &constructfakes.FakeOSValidator{}
		fakeMessenger = &constructfakes.FakeConstructMessenger{}

		vmConstruct = NewVMConstruct(
			context.TODO(),
			fakeRemoteManager,
			"fakeUser",
			"fakePass",
			"fakeVmPath",
			fakeVcenterClient,
			fakeGuestManager,
			fakeWinRMEnabler,
			fakeOSValidator,
			fakeMessenger,
		)

		fakeGuestManager.StartProgramInGuestReturnsOnCall(0, 0, nil)
		fakeGuestManager.ExitCodeForProgramInGuestReturnsOnCall(0, 0, nil)
		versionBuffer := gbytes.NewBuffer()
		_, err := versionBuffer.Write([]byte("dev"))
		Expect(err).NotTo(HaveOccurred())

		fakeGuestManager.DownloadFileInGuestReturns(versionBuffer, 3, nil)
		fakeGuestManager.StartProgramInGuestReturns(0, nil)
	})

	Describe("PrepareVM", func() {
		Context("Validates the OS version of the target machine", func() {
			It("returns failure if the OS Validator returns an error", func() {
				validationError := errors.New("the OS is wrong")
				fakeOSValidator.ValidateReturns(validationError)

				err := vmConstruct.PrepareVM()

				Expect(err).To(MatchError(validationError))
				Expect(fakeVcenterClient.MakeDirectoryCallCount()).To(Equal(0))

				Expect(fakeMessenger.UploadArtifactsStartedCallCount()).To(Equal(0))
			})

			It("prepares the VM if the OS version is correct", func() {
				err := vmConstruct.PrepareVM()

				Expect(err).NotTo(HaveOccurred())
				Expect(fakeMessenger.UploadArtifactsStartedCallCount()).To(Equal(1))
			})
		})

		Context("can create provision directory", func() {
			It("creates it successfully", func() {
				err := vmConstruct.PrepareVM()

				Expect(err).ToNot(HaveOccurred())
				Expect(fakeVcenterClient.MakeDirectoryCallCount()).To(Equal(1))
				Expect(fakeMessenger.CreateProvisionDirStartedCallCount()).To(Equal(1))
				Expect(fakeMessenger.CreateProvisionDirSucceededCallCount()).To(Equal(1))
			})

			It("fails when the provision dir cannot be created", func() {
				mkDirError := errors.New("failed to create dir")
				fakeVcenterClient.MakeDirectoryReturns(mkDirError)

				err := vmConstruct.PrepareVM()

				Expect(fakeVcenterClient.MakeDirectoryCallCount()).To(Equal(1))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to create dir"))
				Expect(fakeMessenger.CreateProvisionDirStartedCallCount()).To(Equal(1))
				Expect(fakeMessenger.CreateProvisionDirSucceededCallCount()).To(Equal(0))
			})
		})

		Context("enable WinRM", func() {
			It("returns failure when it fails to enable winrm", func() {
				execError := errors.New("failed to enable winRM")
				fakeWinRMEnabler.EnableReturns(execError)

				err := vmConstruct.PrepareVM()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to enable winRM"))

				Expect(fakeWinRMEnabler.EnableCallCount()).To(Equal(1))
			})

			It("logs that winrm was successfully enabled", func() {
				err := vmConstruct.PrepareVM()

				Expect(err).NotTo(HaveOccurred())
				Expect(fakeMessenger.EnableWinRMStartedCallCount()).To(Equal(1))
				Expect(fakeMessenger.EnableWinRMSucceededCallCount()).To(Equal(1))
			})
		})

		Context("can connect to VM", func() {
			It("can reach VM and can login to VM", func() {
				err := vmConstruct.PrepareVM()

				Expect(err).To(BeNil())
				Expect(fakeRemoteManager.CanReachVMCallCount()).To(Equal(1))
				Expect(fakeRemoteManager.CanLoginVMCallCount()).To(Equal(1))
			})
			It("returns an error if it cannot reach the VM", func() {
				fakeRemoteManager.CanReachVMReturns(errors.New("can't reach VM"))

				err := vmConstruct.PrepareVM()
				Expect(err).NotTo(BeNil())
				Expect(err).To(MatchError("can't reach VM"))
				Expect(fakeRemoteManager.CanReachVMCallCount()).To(Equal(1))
				Expect(fakeRemoteManager.CanLoginVMCallCount()).To(Equal(0))
			})

			It("should return an error when login fails", func() {
				invalidPwdError := errors.New("login error")
				fakeRemoteManager.CanLoginVMReturns(invalidPwdError)

				err := vmConstruct.PrepareVM()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(invalidPwdError))

				Expect(fakeRemoteManager.CanReachVMCallCount()).To(Equal(1))
				Expect(fakeRemoteManager.CanLoginVMCallCount()).To(Equal(1))
			})

			It("logs that it successfully validated the vm connection", func() {
				err := vmConstruct.PrepareVM()

				Expect(err).NotTo(HaveOccurred())
				Expect(fakeMessenger.ValidateVMConnectionStartedCallCount()).To(Equal(1))
				Expect(fakeMessenger.ValidateVMConnectionSucceededCallCount()).To(Equal(1))
			})

		})

		Context("can upload artifacts", func() {
			Context("Upload all artifacts correctly", func() {
				It("passes successfully", func() {

					err := vmConstruct.PrepareVM()
					Expect(err).ToNot(HaveOccurred())
					vmPath, artifact, dest, user, pass := fakeVcenterClient.UploadArtifactArgsForCall(0)
					Expect(artifact).To(Equal("./LGPO.zip"))
					Expect(vmPath).To(Equal("fakeVmPath"))
					Expect(dest).To(Equal("C:\\provision\\LGPO.zip"))
					Expect(user).To(Equal("fakeUser"))
					Expect(pass).To(Equal("fakePass"))
					Expect(fakeVcenterClient.UploadArtifactCallCount()).To(Equal(2))
					Expect(fakeMessenger.UploadArtifactsStartedCallCount()).To(Equal(1))
					Expect(fakeMessenger.UploadArtifactsSucceededCallCount()).To(Equal(1))

					Expect(fakeMessenger.UploadFileStartedCallCount()).To(Equal(2))
					artifact = fakeMessenger.UploadFileStartedArgsForCall(0)
					Expect(artifact).To(Equal("LGPO"))
					artifact = fakeMessenger.UploadFileStartedArgsForCall(1)
					Expect(artifact).To(Equal("stemcell preparation artifacts"))

					Expect(fakeMessenger.UploadFileSucceededCallCount()).To(Equal(2))
				})

			})

			Context("Fails to upload one or more artifacts", func() {
				It("fails when it cannot upload LGPO", func() {

					uploadError := errors.New("failed to upload LGPO")
					fakeVcenterClient.UploadArtifactReturns(uploadError)

					err := vmConstruct.PrepareVM()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("failed to upload LGPO"))

					vmPath, artifact, _, _, _ := fakeVcenterClient.UploadArtifactArgsForCall(0)
					Expect(artifact).To(Equal("./LGPO.zip"))
					Expect(vmPath).To(Equal("fakeVmPath"))
					Expect(fakeVcenterClient.UploadArtifactCallCount()).To(Equal(1))
					Expect(fakeMessenger.UploadArtifactsStartedCallCount()).To(Equal(1))
					Expect(fakeMessenger.UploadArtifactsSucceededCallCount()).To(Equal(0))
				})

				It("fails when it cannot upload Stemcell Automation scripts", func() {

					uploadError := errors.New("failed to upload stemcell automation")
					fakeVcenterClient.UploadArtifactReturnsOnCall(0, nil)
					fakeVcenterClient.UploadArtifactReturnsOnCall(1, uploadError)

					err := vmConstruct.PrepareVM()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("failed to upload stemcell automation"))

					vmPath, artifact, _, _, _ := fakeVcenterClient.UploadArtifactArgsForCall(0)
					Expect(artifact).To(Equal("./LGPO.zip"))
					Expect(vmPath).To(Equal("fakeVmPath"))
					vmPath, artifact, _, _, _ = fakeVcenterClient.UploadArtifactArgsForCall(1)
					Expect(artifact).To(Equal("./StemcellAutomation.zip"))
					Expect(vmPath).To(Equal("fakeVmPath"))
					Expect(fakeVcenterClient.UploadArtifactCallCount()).To(Equal(2))
					Expect(fakeMessenger.UploadArtifactsStartedCallCount()).To(Equal(1))
					Expect(fakeMessenger.UploadArtifactsSucceededCallCount()).To(Equal(0))
				})
			})
		})

		Context("can extract archives", func() {
			It("returns failure when it fails to extract archive", func() {
				extractError := errors.New("failed to extract archive")
				fakeRemoteManager.ExtractArchiveReturns(extractError)

				err := vmConstruct.PrepareVM()
				Expect(err).To(HaveOccurred())
				Expect(fakeRemoteManager.ExtractArchiveCallCount()).To(Equal(1))
				Expect(err.Error()).To(Equal("failed to extract archive"))
				Expect(fakeMessenger.ExtractArtifactsStartedCallCount()).To(Equal(1))
				Expect(fakeMessenger.ExtractArtifactsSucceededCallCount()).To(Equal(0))
			})

			It("returns success when it properly extracts archive", func() {
				fakeRemoteManager.ExtractArchiveReturns(nil)

				err := vmConstruct.PrepareVM()
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeRemoteManager.ExtractArchiveCallCount()).To(Equal(1))
				source, destination := fakeRemoteManager.ExtractArchiveArgsForCall(0)
				Expect(source).To(Equal("C:\\provision\\StemcellAutomation.zip"))
				Expect(destination).To(Equal("C:\\provision\\"))

				Expect(fakeMessenger.ExtractArtifactsStartedCallCount()).To(Equal(1))
				Expect(fakeMessenger.ExtractArtifactsSucceededCallCount()).To(Equal(1))
			})

		})

		Context("can execute scripts", func() {
			It("returns failure when it fails to execute setup script", func() {
				execError := errors.New("failed to execute setup script")
				fakeRemoteManager.ExecuteCommandReturns(execError)

				err := vmConstruct.PrepareVM()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to execute setup script"))

				Expect(fakeRemoteManager.ExecuteCommandCallCount()).To(Equal(1))
				Expect(fakeMessenger.ExecuteScriptStartedCallCount()).To(Equal(1))
				Expect(fakeMessenger.ExecuteScriptSucceededCallCount()).To(Equal(0))
			})

			It("returns success when it properly executes the setup script", func() {

				err := vmConstruct.PrepareVM()
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeRemoteManager.ExecuteCommandCallCount()).To(Equal(1))
				command := fakeRemoteManager.ExecuteCommandArgsForCall(0)
				Expect(command).To(Equal("powershell.exe C:\\provision\\Setup.ps1"))

				Expect(fakeMessenger.ExecuteScriptStartedCallCount()).To(Equal(1))
				Expect(fakeMessenger.ExecuteScriptSucceededCallCount()).To(Equal(1))
			})

		})
	})
})
