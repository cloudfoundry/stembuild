package construct_test

import (
	"errors"
	. "github.com/cloudfoundry-incubator/stembuild/construct"
	"github.com/cloudfoundry-incubator/stembuild/construct/constructfakes"
	"github.com/cloudfoundry-incubator/stembuild/remotemanager/remotemanagerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("construct_helpers", func() {
	var (
		fakeRemoteManager *remotemanagerfakes.FakeRemoteManager
		mockVMConstruct   *VMConstruct
		fakeVcenterClient *constructfakes.FakeIaasClient
	)

	BeforeEach(func() {
		fakeRemoteManager = &remotemanagerfakes.FakeRemoteManager{}
		fakeVcenterClient = &constructfakes.FakeIaasClient{}
		mockVMConstruct = NewMockVMConstruct(fakeRemoteManager, fakeVcenterClient, "fakeVmPath", "fakeUser", "fakePass")
	})

	Describe("CanConnectToVM", func() {
		It("should not return an error if vm & credential are valid", func() {
			fakeRemoteManager.CanReachVMReturns(nil)
			fakeRemoteManager.CanLoginVMReturns(nil)

			err := mockVMConstruct.CanConnectToVM()
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeRemoteManager.CanReachVMCallCount()).To(Equal(1))
			Expect(fakeRemoteManager.CanLoginVMCallCount()).To(Equal(1))

		})

		It("should return an error if vm is invalid", func() {
			invalidVMError := errors.New("invalid vm")
			fakeRemoteManager.CanReachVMReturns(invalidVMError)

			err := mockVMConstruct.CanConnectToVM()
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(invalidVMError))

			Expect(fakeRemoteManager.CanReachVMCallCount()).To(Equal(1))
			Expect(fakeRemoteManager.CanLoginVMCallCount()).To(Equal(0))
		})

		It("should return an error if username/password is invalid", func() {
			invalidPwdError := errors.New("invalid password")
			fakeRemoteManager.CanReachVMReturns(nil)
			fakeRemoteManager.CanLoginVMReturns(invalidPwdError)

			err := mockVMConstruct.CanConnectToVM()
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(invalidPwdError))

			Expect(fakeRemoteManager.CanReachVMCallCount()).To(Equal(1))
			Expect(fakeRemoteManager.CanLoginVMCallCount()).To(Equal(1))

		})
	})

	Describe("UploadArtifact", func() {

		Context("Upload all artifacts correctly", func() {
			It("passes successfully", func() {

				fakeVcenterClient.MakeDirectoryReturns(nil)
				fakeVcenterClient.UploadArtifactReturns(nil)
				err := mockVMConstruct.UploadArtifact()
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeVcenterClient.MakeDirectoryCallCount()).To(Equal(1))
				vmPath, artifact, dest, user, pass := fakeVcenterClient.UploadArtifactArgsForCall(0)
				Expect(artifact).To(Equal("./LGPO.zip"))
				Expect(vmPath).To(Equal("fakeVmPath"))
				Expect(dest).To(Equal("C:\\provision\\LGPO.zip"))
				Expect(user).To(Equal("fakeUser"))
				Expect(pass).To(Equal("fakePass"))
				Expect(fakeVcenterClient.UploadArtifactCallCount()).To(Equal(2))
			})

		})

		Context("Fails to upload one or more artifacts", func() {
			It("fails when it cannot upload LGPO", func() {

				uploadError := errors.New("failed to upload LGPO")
				fakeVcenterClient.UploadArtifactReturns(uploadError)

				err := mockVMConstruct.UploadArtifact()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to upload LGPO"))

				vmPath, artifact, _, _, _ := fakeVcenterClient.UploadArtifactArgsForCall(0)
				Expect(artifact).To(Equal("./LGPO.zip"))
				Expect(vmPath).To(Equal("fakeVmPath"))
				Expect(fakeVcenterClient.UploadArtifactCallCount()).To(Equal(1))
			})

			It("fails when it cannot upload Stemcell Automation scripts", func() {

				uploadError := errors.New("failed to upload stemcell automation")
				fakeVcenterClient.UploadArtifactReturnsOnCall(0, nil)
				fakeVcenterClient.UploadArtifactReturnsOnCall(1, uploadError)

				err := mockVMConstruct.UploadArtifact()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to upload stemcell automation"))

				vmPath, artifact, _, _, _ := fakeVcenterClient.UploadArtifactArgsForCall(0)
				Expect(artifact).To(Equal("./LGPO.zip"))
				Expect(vmPath).To(Equal("fakeVmPath"))
				vmPath, artifact, _, _, _ = fakeVcenterClient.UploadArtifactArgsForCall(1)
				Expect(artifact).To(Equal("./StemcellAutomation.zip"))
				Expect(vmPath).To(Equal("fakeVmPath"))
				Expect(fakeVcenterClient.UploadArtifactCallCount()).To(Equal(2))
			})

			It("fails when the provision dir cannot be created", func() {

				mkDirError := errors.New("failed to create dir")
				fakeVcenterClient.MakeDirectoryReturns(mkDirError)

				err := mockVMConstruct.UploadArtifact()
				Expect(fakeVcenterClient.MakeDirectoryCallCount()).To(Equal(1))

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to create dir"))

				Expect(fakeVcenterClient.UploadArtifactCallCount()).To(Equal(0))
			})
		})
	})

	Describe("ExtractArchive", func() {

		It("returns failure when it fails to extract archive", func() {
			extractError := errors.New("failed to extract archive")
			fakeRemoteManager.ExtractArchiveReturns(extractError)

			err := mockVMConstruct.ExtractArchive()
			Expect(err).To(HaveOccurred())
			Expect(fakeRemoteManager.ExtractArchiveCallCount()).To(Equal(1))
			Expect(err.Error()).To(Equal("failed to extract archive"))
		})

		It("returns success when it properly extracts archive", func() {
			fakeRemoteManager.ExtractArchiveReturns(nil)

			err := mockVMConstruct.ExtractArchive()
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeRemoteManager.ExtractArchiveCallCount()).To(Equal(1))

		})
	})

	Describe("ExecuteSetupScript", func() {
		It("returns failure when it fails to execute setup script", func() {
			execError := errors.New("failed to execute setup script")
			fakeRemoteManager.ExecuteCommandReturns(execError)

			err := mockVMConstruct.ExecuteSetupScript()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to execute setup script"))

			Expect(fakeRemoteManager.ExecuteCommandCallCount()).To(Equal(1))
		})

		It("returns success when it properly executes the setup script", func() {
			fakeRemoteManager.ExecuteCommandReturns(nil)

			err := mockVMConstruct.ExecuteSetupScript()
			Expect(err).ToNot(HaveOccurred())

			Expect(fakeRemoteManager.ExecuteCommandCallCount()).To(Equal(1))

		})
	})
})
