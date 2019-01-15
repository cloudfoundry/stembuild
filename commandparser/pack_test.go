package commandparser_test

import (
	"errors"
	"flag"
	. "github.com/cloudfoundry-incubator/stembuild/commandparser"
	. "github.com/cloudfoundry-incubator/stembuild/filesystem/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var _ = Describe("pack", func() {
	// Focus of this test is not to test the Flags.Parse functionality as much
	// as to test that the command line flags values are stored in the expected
	// struct variables. This adds a bit of protection when renaming flag parameters.
	Describe("SetFlags", func() {

		var f *flag.FlagSet
		var PkgCmd *PackageCmd

		BeforeEach(func() {
			f = flag.NewFlagSet("test", flag.ExitOnError)
			PkgCmd = &PackageCmd{}
			PkgCmd.SetFlags(f)
		})

		var longformArgs = []string{"-vmdk", "some_vmdk_file",
			"-os", "1803",
			"-stemcell-version", "1803.45",
			"-outputDir", "some_output_dir",
		}
		var shortformArgs = []string{"-vmdk", "some_vmdk_file",
			"-os", "1803",
			"-s", "1803.45",
			"-o", "some_output_dir",
		}

		Context("a vmdk file is specified as a flag parameter", func() {
			It("then the vmdk file name is stored", func() {
				err := f.Parse(longformArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(PkgCmd.GetVMDK()).To(Equal("some_vmdk_file"))
			})
		})

		Context("a os stemcellVersion is specified as a flag parameter", func() {
			It("then the os stemcellVersion is stored", func() {
				err := f.Parse(longformArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(PkgCmd.GetOS()).To(Equal("1803"))
			})
		})

		Context("a stemcell stemcellVersion is specified as a flag parameter", func() {
			It("when using the long form the stemcell stemcellVersion is stored", func() {
				err := f.Parse(longformArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(PkgCmd.GetStemcellVersion()).To(Equal("1803.45"))
			})

			It("when using the short form the stemcell stemcellVersion is stored", func() {
				err := f.Parse(shortformArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(PkgCmd.GetStemcellVersion()).To(Equal("1803.45"))
			})

		})

		Context("an output directory is specified as a flag parameter", func() {
			It("when using the long form the directory is stored", func() {
				err := f.Parse(longformArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(PkgCmd.GetOutputDir()).To(Equal("some_output_dir"))
			})

			It("when using the short form the directory is stored", func() {
				err := f.Parse(shortformArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(PkgCmd.GetOutputDir()).To(Equal("some_output_dir"))
			})

		})

	})

	Describe("validateFreeSpaceForPackage", func() {
		var (
			mockCtrl       *gomock.Controller
			mockFileSystem *MockFileSystem
		)
		const gigFreeSpace = uint64(1024 * 1024 * 1024)
		const lowFreeSpace = uint64(20)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockFileSystem = NewMockFileSystem(mockCtrl)
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		Context("There is enough space on disk", func() {
			It("returns true", func() {
				mockFileSystem.EXPECT().GetAvailableDiskSpace(gomock.Any()).Return(gigFreeSpace, nil).AnyTimes()

				pkgCmd := PackageCmd{}
				pkgCmd.SetVMDK(filepath.Join("..", "test", "data", "expected.vmdk"))

				validSpace, _, err := pkgCmd.ValidateFreeSpaceForPackage(mockFileSystem)
				Expect(err).ToNot(HaveOccurred())
				Expect(validSpace).To(BeTrue())
			})
		})

		Context("There is not enough space on disk", func() {
			It("returns false", func() {
				mockFileSystem.EXPECT().GetAvailableDiskSpace(gomock.Any()).Return(lowFreeSpace, nil).AnyTimes()

				pkgCmd := PackageCmd{}
				pkgCmd.SetVMDK(filepath.Join("..", "test", "data", "expected.vmdk"))

				validSpace, spaceNeeded, err := pkgCmd.ValidateFreeSpaceForPackage(mockFileSystem)
				Expect(err).ToNot(HaveOccurred())
				Expect(validSpace).To(BeFalse())
				Expect(spaceNeeded).To(Equal(uint64(563085292)))
			})
		})

		Context("Returns an error", func() {
			It("when the vmdk doesn't exist", func() {
				mockFileSystem.EXPECT().GetAvailableDiskSpace(gomock.Any()).Return(gigFreeSpace, nil).AnyTimes()

				pkgCmd := PackageCmd{}
				pkgCmd.SetVMDK(filepath.Join("..", "test", "data", "nonexistent.vmdk"))

				validSpace, _, err := pkgCmd.ValidateFreeSpaceForPackage(mockFileSystem)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(HavePrefix("could not get vmdk info: "))
				Expect(validSpace).To(BeFalse())
			})

			It("when it fails to check the available disk space", func() {
				mockFileSystem.EXPECT().GetAvailableDiskSpace(gomock.Any()).Return(gigFreeSpace, errors.New("some check error")).AnyTimes()

				pkgCmd := PackageCmd{}
				pkgCmd.SetVMDK(filepath.Join("..", "test", "data", "expected.vmdk"))

				validSpace, _, err := pkgCmd.ValidateFreeSpaceForPackage(mockFileSystem)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("could not check free space on disk: some check error"))
				Expect(validSpace).To(BeFalse())
			})
		})
	})
})
