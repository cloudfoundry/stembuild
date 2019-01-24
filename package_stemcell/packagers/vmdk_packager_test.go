package packagers_test

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang/mock/gomock"

	. "github.com/cloudfoundry-incubator/stembuild/filesystem/mock"
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/package_parameters"
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/packagers"
	"github.com/cloudfoundry-incubator/stembuild/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VmdkPackager", func() {
	var tmpDir string
	var stembuildConfig package_parameters.VmdkPackageParameters
	var c packagers.VmdkPackager

	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		stembuildConfig = package_parameters.VmdkPackageParameters{
			OSVersion: "2012R2",
			Version:   "1200.1",
		}

		c = packagers.VmdkPackager{
			Stop:         make(chan struct{}),
			Debugf:       func(format string, a ...interface{}) {},
			BuildOptions: stembuildConfig,
		}
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	Describe("vmdk", func() {
		Context("valid vmdk file specified", func() {
			It("should be valid", func() {

				vmdk, err := ioutil.TempFile("", "temp.vmdk")
				Expect(err).ToNot(HaveOccurred())
				defer os.Remove(vmdk.Name())

				valid, err := packagers.IsValidVMDK(vmdk.Name())
				Expect(err).To(BeNil())
				Expect(valid).To(BeTrue())
			})
		})
		Context("invalid vmdk file specified", func() {
			It("should be invalid", func() {
				valid, err := packagers.IsValidVMDK(filepath.Join("..", "out", "invalid"))
				Expect(err).To(HaveOccurred())
				Expect(valid).To(BeFalse())
			})
		})
	})

	Describe("CreateImage", func() {

		It("successfully creates an image tarball", func() {
			c.BuildOptions.VMDKFile = filepath.Join("..", "..", "test", "data", "expected.vmdk")
			err := c.CreateImage()
			Expect(err).NotTo(HaveOccurred())

			// the image will be saved to the VmdkPackager's temp directory
			tmpdir, err := c.TempDir()
			Expect(err).NotTo(HaveOccurred())

			outputImagePath := filepath.Join(tmpdir, "image")
			Expect(c.Image).To(Equal(outputImagePath))

			// Make sure the sha1 sum is correct
			h := sha1.New()
			f, err := os.Open(c.Image)
			Expect(err).NotTo(HaveOccurred())

			_, err = io.Copy(h, f)
			Expect(err).NotTo(HaveOccurred())

			actualShasum := fmt.Sprintf("%x", h.Sum(nil))
			Expect(c.Sha1sum).To(Equal(actualShasum))

			// expect the image ova to contain only the following file names
			expectedNames := []string{
				"image.ovf",
				"image.mf",
				"image-disk1.vmdk",
			}

			imageDir, err := helpers.ExtractGzipArchive(c.Image)
			Expect(err).NotTo(HaveOccurred())
			list, err := ioutil.ReadDir(imageDir)
			Expect(err).NotTo(HaveOccurred())

			var names []string
			infos := make(map[string]os.FileInfo)
			for _, fi := range list {
				names = append(names, fi.Name())
				infos[fi.Name()] = fi
			}
			Expect(names).To(ConsistOf(expectedNames))

			// the vmx template should generate an ovf file that
			// does not contain an ethernet section.
			ovf := filepath.Join(imageDir, "image.ovf")
			ovfFile, err := helpers.ReadFile(ovf)
			Expect(err).NotTo(HaveOccurred())
			Expect(ovfFile).NotTo(MatchRegexp(`(?i)ethernet`))
		})
	})

	Describe("ValidateFreeSpaceForPackage", func() {
		var (
			mockCtrl       *gomock.Controller
			mockFileSystem *MockFileSystem
		)

		Context("When VMDK file is invalid", func() {
			It("returns an error", func() {
				c.BuildOptions.VMDKFile = ""

				mockCtrl = gomock.NewController(GinkgoT())
				mockFileSystem = NewMockFileSystem(mockCtrl)

				err := c.ValidateFreeSpaceForPackage(mockFileSystem)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("could not get vmdk info"))

			})
		})

		Context("When filesystem has enough free space for stemcell (twice the size of the expected free space)", func() {
			It("does not return an error", func() {
				c.BuildOptions.VMDKFile = filepath.Join("..", "..", "test", "data", "expected.vmdk")

				mockCtrl = gomock.NewController(GinkgoT())
				mockFileSystem = NewMockFileSystem(mockCtrl)

				vmdkFile, err := os.Stat(c.BuildOptions.VMDKFile)
				Expect(err).ToNot(HaveOccurred())

				testVmdkSize := vmdkFile.Size()
				expectFreeSpace := uint64(testVmdkSize)*2 + (packagers.Gigabyte / 2)

				directoryPath := filepath.Dir(c.BuildOptions.VMDKFile)
				mockFileSystem.EXPECT().GetAvailableDiskSpace(directoryPath).Return(uint64(expectFreeSpace*2), nil).AnyTimes()

				err = c.ValidateFreeSpaceForPackage(mockFileSystem)
				Expect(err).To(Not(HaveOccurred()))

			})
		})
		Context("When filesystem does not have enough free space for stemcell (half the size of the expected free space", func() {
			It("returns error", func() {
				c.BuildOptions.VMDKFile = filepath.Join("..", "..", "test", "data", "expected.vmdk")

				mockCtrl = gomock.NewController(GinkgoT())
				mockFileSystem = NewMockFileSystem(mockCtrl)

				vmdkFile, err := os.Stat(c.BuildOptions.VMDKFile)
				Expect(err).ToNot(HaveOccurred())

				testVmdkSize := vmdkFile.Size()
				expectFreeSpace := uint64(testVmdkSize)*2 + (packagers.Gigabyte / 2)

				directoryPath := filepath.Dir(c.BuildOptions.VMDKFile)
				mockFileSystem.EXPECT().GetAvailableDiskSpace(directoryPath).Return(uint64(expectFreeSpace/2), nil).AnyTimes()

				err = c.ValidateFreeSpaceForPackage(mockFileSystem)

				Expect(err).To(HaveOccurred())

				expectedErrorMsg := fmt.Sprintf("Not enough space to create stemcell. Free up ")
				Expect(err.Error()).To(ContainSubstring(expectedErrorMsg))
			})
		})

		Context("When filesystem fails to provide free space", func() {
			It("returns error specifying that given disk could not provide free space", func() {
				c.BuildOptions.VMDKFile = filepath.Join("..", "..", "test", "data", "expected.vmdk")

				mockCtrl = gomock.NewController(GinkgoT())
				mockFileSystem = NewMockFileSystem(mockCtrl)

				directoryPath := filepath.Dir(c.BuildOptions.VMDKFile)
				mockFileSystem.EXPECT().GetAvailableDiskSpace(directoryPath).Return(uint64(4), errors.New("some error")).AnyTimes()

				err := c.ValidateFreeSpaceForPackage(mockFileSystem)

				Expect(err).To(HaveOccurred())
				expectedErrorMsg := fmt.Sprintf("could not check free space on disk: ")
				Expect(err.Error()).To(ContainSubstring(expectedErrorMsg))
			})
		})
	})
})
