package construct_test

import (
	"errors"
	. "github.com/cloudfoundry-incubator/stembuild/construct"
	. "github.com/cloudfoundry-incubator/stembuild/remotemanager/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("construct_helpers", func() {
	var (
		mockCtrl          *gomock.Controller
		mockRemoteManager *MockRemoteManager
		mockVMConstruct   *VMConstruct
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRemoteManager = NewMockRemoteManager(mockCtrl)
		mockVMConstruct = NewMockVMConstruct(mockRemoteManager)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("UploadArtifact", func() {

		Context("Upload all artifacts correctly", func() {
			It("passes successfully", func() {
				mockRemoteManager.EXPECT().UploadArtifact(gomock.Any(), gomock.Any()).Return(nil).Times(2)

				err := mockVMConstruct.UploadArtifact()
				Expect(err).ToNot(HaveOccurred())
			})

		})

		Context("Fails to upload one or more artifacts", func() {
			It("fails when it cannot upload LGPO", func() {
				mockRemoteManager.EXPECT().UploadArtifact("./LGPO.zip", gomock.Any()).Return(errors.New("failed to upload LGPO")).Times(1)
				mockRemoteManager.EXPECT().UploadArtifact("./StemcellAutomation.zip", gomock.Any()).Return(nil).Times(0)

				err := mockVMConstruct.UploadArtifact()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to upload LGPO"))
			})

			It("fails when it cannot upload Stemcell Automation scripts", func() {
				mockRemoteManager.EXPECT().UploadArtifact("./LGPO.zip", gomock.Any()).Return(nil).Times(1)
				mockRemoteManager.EXPECT().UploadArtifact("./StemcellAutomation.zip", gomock.Any()).Return(errors.New("failed to upload stemcell automation")).Times(1)

				err := mockVMConstruct.UploadArtifact()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to upload stemcell automation"))
			})
		})

	})

	Describe("ExtractArchive", func() {

		It("returns failure when it fails to extract archive", func() {
			mockRemoteManager.EXPECT().ExtractArchive(gomock.Any(), gomock.Any()).Return(errors.New("failed to extract archive")).Times(1)

			err := mockVMConstruct.ExtractArchive()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to extract archive"))
		})

		It("returns success when it properly extrects archive", func() {
			mockRemoteManager.EXPECT().ExtractArchive(gomock.Any(), gomock.Any()).Return(nil).Times(1)

			err := mockVMConstruct.ExtractArchive()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("ExecuteSetupScript", func() {
		It("returns failure when it fails to execute setup script", func() {
			mockRemoteManager.EXPECT().ExecuteCommand(gomock.Any()).Return(errors.New("failed to execute setup script")).Times(1)

			err := mockVMConstruct.ExecuteSetupScript()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to execute setup script"))
		})

		It("returns success when it properly executes the setup script", func() {
			mockRemoteManager.EXPECT().ExecuteCommand(gomock.Any()).Return(nil).Times(1)

			err := mockVMConstruct.ExecuteSetupScript()
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
