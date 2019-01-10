package commandparser_test

import (
	"errors"

	. "github.com/cloudfoundry-incubator/stembuild/commandparser"
	. "github.com/cloudfoundry-incubator/stembuild/filesystem/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("inputs", func() {

	Describe("HasAvailableDiskSpace", func() {
		var (
			mockCtrl       *gomock.Controller
			mockFileSystem *MockFileSystem
		)

		It("Has enough free space", func() {
			mockCtrl = gomock.NewController(GinkgoT())
			defer mockCtrl.Finish()
			mockFileSystem = NewMockFileSystem(mockCtrl)

			mockFileSystem.EXPECT().GetAvailableDiskSpace("/").Return(uint64(8), nil).AnyTimes()

			hasSpace, _, err := HasAtLeastFreeDiskSpace(4, mockFileSystem, "/")
			Expect(err).To(Not(HaveOccurred()))
			Expect(hasSpace).To(BeTrue())
		})

		It("Not enough free space", func() {
			mockCtrl = gomock.NewController(GinkgoT())
			defer mockCtrl.Finish()
			mockFileSystem = NewMockFileSystem(mockCtrl)

			mockFileSystem.EXPECT().GetAvailableDiskSpace("/").Return(uint64(4), nil).AnyTimes()

			hasSpace, requiredSpace, err := HasAtLeastFreeDiskSpace(8, mockFileSystem, "/")
			Expect(err).To(Not(HaveOccurred()))
			Expect(hasSpace).To(BeFalse())
			Expect(requiredSpace).To(Equal(uint64(4)))
		})

		It("fails on error", func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockFileSystem = NewMockFileSystem(mockCtrl)

			mockFileSystem.EXPECT().GetAvailableDiskSpace("/").Return(uint64(4), errors.New("some error")).AnyTimes()

			hasSpace, _, err := HasAtLeastFreeDiskSpace(8, mockFileSystem, "/")
			Expect(err).To(HaveOccurred())
			Expect(hasSpace).To(BeFalse())
		})
	})
})
