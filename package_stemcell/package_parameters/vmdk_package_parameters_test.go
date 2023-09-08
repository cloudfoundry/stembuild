package package_parameters_test

import (
	"fmt"
	"math/rand"

	"github.com/cloudfoundry/stembuild/package_stemcell/package_parameters"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VmdkPackageParameters", func() {
	Context("CopyFrom", func() {
		var (
			dest package_parameters.VmdkPackageParameters
			src  package_parameters.VmdkPackageParameters
		)

		BeforeEach(func() {
			src = package_parameters.VmdkPackageParameters{}
			dest = package_parameters.VmdkPackageParameters{}
		})

		JustBeforeEach(func() {
			dest.CopyFrom(src)
		})

		Context("OSVersion", func() {
			Context("when src specifies an OSVersion and dest does not", func() {
				BeforeEach(func() {
					src.OSVersion = fmt.Sprintf("banana%d.%d", rand.Intn(2000), rand.Intn(2000))
				})

				It("copies src.OSVersion into dest.OSVersion", func() {
					Expect(dest.OSVersion).To(Equal(src.OSVersion))
				})
			})

			Context("when src specifies an OSVersion and dest specifies an OSVersion", func() {
				var expectedDestOSVersion string

				BeforeEach(func() {
					src.OSVersion = fmt.Sprintf("banana%d.%d", rand.Intn(2000), rand.Intn(2000))
					dest.OSVersion = fmt.Sprintf("banana%d.%d", rand.Intn(2000), rand.Intn(2000))
					expectedDestOSVersion = dest.OSVersion
				})

				It("retains dest.OSVersion's original value", func() {
					Expect(dest.OSVersion).To(Equal(expectedDestOSVersion))
				})
			})

			Context("when src does not specify an OSVersion and dest does not specify an OSVersion", func() {
				It("should do nothing", func() {
					Expect(dest.OSVersion).To(BeEmpty())
				})
			})

			Context("when dest does specify a OSVersion and src does not specify a OSVersion", func() {
				var expectedDestOSVersion string

				BeforeEach(func() {
					dest.OSVersion = fmt.Sprintf("banana%d.%d", rand.Intn(2000), rand.Intn(2000))
					expectedDestOSVersion = dest.OSVersion
				})
				It("should do nothing", func() {
					Expect(dest.OSVersion).To(Equal(expectedDestOSVersion))
				})
			})
		})

		Context("OutputDir", func() {
			Context("when src specifies an OutputDir and dest does not", func() {
				BeforeEach(func() {
					src.OutputDir = fmt.Sprintf("foo%d/%d", rand.Intn(2000), rand.Intn(2000))
				})

				It("does nothing", func() {
					Expect(dest.OutputDir).To(Equal(""))
				})
			})

			Context("when src specifies an OutputDir and dest specifies an OutputDir", func() {
				var expectedDestOutputDir string

				BeforeEach(func() {
					src.OutputDir = fmt.Sprintf("foo%d/%d", rand.Intn(2000), rand.Intn(2000))
					dest.OutputDir = fmt.Sprintf("foo%d/%d", rand.Intn(2000), rand.Intn(2000))
					expectedDestOutputDir = dest.OutputDir
				})

				It("retains dest.OutputDir's original value", func() {
					Expect(dest.OutputDir).To(Equal(expectedDestOutputDir))
				})
			})

			Context("when src does not specify an OutputDir and dest does not specify an OutputDir", func() {
				It("should do nothing", func() {
					Expect(dest.OutputDir).To(BeEmpty())
				})
			})

			Context("when dest does specify a OutputDir and src does not specify a OutputDir", func() {
				var expectedDestOutputDir string

				BeforeEach(func() {
					dest.OutputDir = fmt.Sprintf("foo%d/%d", rand.Intn(2000), rand.Intn(2000))
					expectedDestOutputDir = dest.OutputDir
				})
				It("should do nothing", func() {
					Expect(dest.OutputDir).To(Equal(expectedDestOutputDir))
				})
			})
		})

		Context("Version", func() {
			Context("when src specifies a version and dest does not", func() {
				BeforeEach(func() {
					src.Version = fmt.Sprintf("%d.%d", rand.Intn(2000), rand.Intn(2000))
				})

				It("copies src.Version into dest.Version", func() {
					Expect(dest.Version).To(Equal(src.Version))
				})
			})

			Context("when src specifies a version and dest specifies a version", func() {
				var expectedDestVersion string

				BeforeEach(func() {
					src.Version = fmt.Sprintf("%d.%d", rand.Intn(2000), rand.Intn(2000))
					dest.Version = fmt.Sprintf("%d.%d", rand.Intn(2000), rand.Intn(2000))
					expectedDestVersion = dest.Version
				})

				It("retains dest.Version's original value", func() {
					Expect(dest.Version).To(Equal(expectedDestVersion))
				})
			})

			Context("when src does not specify a version and dest does not specify a version", func() {
				It("should do nothing", func() {
					Expect(dest.Version).To(BeEmpty())
				})
			})

			Context("when dest does specify a Version and src does not specify a Version", func() {
				var expectedDestVersion string

				BeforeEach(func() {
					dest.Version = fmt.Sprintf("%d.%d", rand.Intn(2000), rand.Intn(2000))
					expectedDestVersion = dest.Version
				})
				It("should do nothing", func() {
					Expect(dest.Version).To(Equal(expectedDestVersion))
				})
			})
		})

		Context("VMDKFile", func() {
			Context("when src specifies an VMDKFile and dest does not", func() {
				BeforeEach(func() {
					src.VMDKFile = fmt.Sprintf("orange%d.vmdk", rand.Intn(2000))
				})

				It("copies src.VMDKFile into dest.VMDKFile", func() {
					Expect(dest.VMDKFile).To(Equal(src.VMDKFile))
				})
			})

			Context("when src specifies a VMDKFile and dest specifies a VMDKFile", func() {
				var expectedDestVMDKFile string

				BeforeEach(func() {
					src.VMDKFile = fmt.Sprintf("orange%d.vmdk", rand.Intn(2000))
					dest.VMDKFile = fmt.Sprintf("orange%d.vmdk", rand.Intn(2000))
					expectedDestVMDKFile = dest.VMDKFile
				})

				It("retains dest.VMDKFile's original value", func() {
					Expect(dest.VMDKFile).To(Equal(expectedDestVMDKFile))
				})
			})

			Context("when src does not specify a VMDKFile and dest does not specify a VMDKFile", func() {
				It("should do nothing", func() {
					Expect(dest.VMDKFile).To(BeEmpty())
				})
			})

			Context("when dest does specify a VMDKFile and src does not specify a VMDKFile", func() {
				var expectedDestVMDKFile string

				BeforeEach(func() {
					expectedDestVMDKFile = dest.VMDKFile
				})
				It("should do nothing", func() {
					Expect(dest.VMDKFile).To(Equal(expectedDestVMDKFile))
				})
			})
		})

		Context("Multiple fields", func() {
			Context("when some fields are set in src and another, somewhat overlapping, set of fields is set in dest", func() {
				BeforeEach(func() {
					dest.OutputDir = "needful"
					dest.Version = "do"
					dest.VMDKFile = "not"
					src.OSVersion = "the"
					src.OutputDir = "bear"
				})

				It("copies into dest only those fields which are empty in dest", func() {
					expected := package_parameters.VmdkPackageParameters{
						OSVersion: "the",
						OutputDir: "needful",
						Version:   "do",
						VMDKFile:  "not",
					}
					Expect(dest).To(Equal(expected))
				})
			})
		})
	})

})
