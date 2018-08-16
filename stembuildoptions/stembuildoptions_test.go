package stembuildoptions_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/stembuild/stembuildoptions"
)

var _ = Describe("StembuildOptions", func() {
	Context("CopyFrom", func() {
		var (
			dest StembuildOptions
			src  StembuildOptions
		)

		BeforeEach(func() {
			src = StembuildOptions{}
			dest = StembuildOptions{}
		})

		JustBeforeEach(func() {
			dest.CopyFrom(src)
		})

		Context("PatchFile", func() {
			Context("when src specifies an PatchFile and dest does not", func() {
				BeforeEach(func() {
					src.PatchFile = fmt.Sprintf("blackberry%d", rand.Intn(2000))
				})

				It("copies src.PatchFile into dest.PatchFile", func() {
					Expect(dest.PatchFile).To(Equal(src.PatchFile))
				})
			})

			Context("when src specifies a PatchFile and dest specifies a PatchFile", func() {
				var expectedDestPatchFile string

				BeforeEach(func() {
					src.PatchFile = fmt.Sprintf("blackberry%d", rand.Intn(2000))
					dest.PatchFile = fmt.Sprintf("blackberry%d", rand.Intn(2000))
					expectedDestPatchFile = dest.PatchFile
				})

				It("retains dest.PatchFile's original value", func() {
					Expect(dest.PatchFile).To(Equal(expectedDestPatchFile))
				})
			})

			Context("when src does not specify a PatchFile and dest does not specify a PatchFile", func() {
				It("should do nothing", func() {
					Expect(dest.PatchFile).To(BeEmpty())
				})
			})

			Context("when dest does specify a PatchFile and src does not specify a PatchFile", func() {
				var expectedDestPatchfile string

				BeforeEach(func() {
					dest.PatchFile = fmt.Sprintf("blackberry%d", rand.Intn(2000))
					expectedDestPatchfile = dest.PatchFile
				})
				It("should do nothing", func() {
					Expect(dest.PatchFile).To(Equal(expectedDestPatchfile))
				})
			})
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

		Context("VHDFile", func() {
			Context("when src specifies an VHDFile and dest does not", func() {
				BeforeEach(func() {
					src.VHDFile = fmt.Sprintf("bar%d.vhd", rand.Intn(2000))
				})

				It("copies src.VHDFile into dest.VHDFile", func() {
					Expect(dest.VHDFile).To(Equal(src.VHDFile))
				})
			})

			Context("when src specifies a VHDFile and dest specifies a VHDFile", func() {
				var expectedDestVHDFile string

				BeforeEach(func() {
					src.VHDFile = fmt.Sprintf("bar%d.vhd", rand.Intn(2000))
					dest.VHDFile = fmt.Sprintf("bar%d.vhd", rand.Intn(2000))
					expectedDestVHDFile = dest.VHDFile
				})

				It("retains dest.VHDFile's original value", func() {
					Expect(dest.VHDFile).To(Equal(expectedDestVHDFile))
				})
			})

			Context("when src does not specify a VHDFile and dest does not specify a VHDFile", func() {
				It("should do nothing", func() {
					Expect(dest.VHDFile).To(BeEmpty())
				})
			})

			Context("when dest does specify a VHDFile and src does not specify a VHDFile", func() {
				var expectedDestVHDFile string

				BeforeEach(func() {
					dest.VHDFile = fmt.Sprintf("bar%d.vhd", rand.Intn(2000))
					expectedDestVHDFile = dest.VHDFile
				})
				It("should do nothing", func() {
					Expect(dest.VHDFile).To(Equal(expectedDestVHDFile))
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
					dest.VHDFile = fmt.Sprintf("orange%d.vmdk", rand.Intn(2000))
					expectedDestVMDKFile = dest.VMDKFile
				})
				It("should do nothing", func() {
					Expect(dest.VMDKFile).To(Equal(expectedDestVMDKFile))
				})
			})
		})

		Context("PatchFileChecksum", func() {
			Context("when src specifies an PatchFileChecksum and dest does not", func() {
				BeforeEach(func() {
					src.PatchFileChecksum = fmt.Sprintf("%d%d%d", rand.Intn(2000), rand.Intn(2000), rand.Intn(2000))
				})

				It("copies src.PatchFileChecksum into dest.PatchFileChecksum", func() {
					Expect(dest.PatchFileChecksum).To(Equal(src.PatchFileChecksum))
				})
			})

			Context("when src specifies a PatchFileChecksum and dest specifies a PatchFileChecksum", func() {
				var expectedDestPatchFileChecksum string

				BeforeEach(func() {
					src.PatchFileChecksum = fmt.Sprintf("%d%d%d", rand.Intn(2000), rand.Intn(2000), rand.Intn(2000))
					dest.PatchFileChecksum = fmt.Sprintf("%d%d%d", rand.Intn(2000), rand.Intn(2000), rand.Intn(2000))
					expectedDestPatchFileChecksum = dest.PatchFileChecksum
				})

				It("retains dest.PatchFileChecksum's original value", func() {
					Expect(dest.PatchFileChecksum).To(Equal(expectedDestPatchFileChecksum))
				})
			})

			Context("when src does not specify a PatchFileChecksum and dest does not specify a PatchFileChecksum", func() {
				It("should do nothing", func() {
					Expect(dest.PatchFileChecksum).To(BeEmpty())
				})
			})

			Context("when dest does specify a PatchFileChecksum and src does not specify a PatchFileChecksum", func() {
				var expectedDestPatchFileChecksum string

				BeforeEach(func() {
					dest.PatchFileChecksum = fmt.Sprintf("orange%d.vmdk", rand.Intn(2000))
					expectedDestPatchFileChecksum = dest.PatchFileChecksum
				})
				It("should do nothing", func() {
					Expect(dest.PatchFileChecksum).To(Equal(expectedDestPatchFileChecksum))
				})
			})
		})

		Context("VHDFileChecksum", func() {
			Context("when src specifies an PatchFileChecksum and dest does not", func() {
				BeforeEach(func() {
					src.VHDFileChecksum = fmt.Sprintf("%d%d%d", rand.Intn(2000), rand.Intn(2000), rand.Intn(2000))
				})

				It("copies src.VHDFileChecksum into dest.VHDFileChecksum", func() {
					Expect(dest.VHDFileChecksum).To(Equal(src.VHDFileChecksum))
				})
			})

			Context("when src specifies a VHDFileChecksum and dest specifies a VHDFileChecksum", func() {
				var expectedDestVHDFileChecksum string

				BeforeEach(func() {
					src.VHDFileChecksum = fmt.Sprintf("%d%d%d", rand.Intn(2000), rand.Intn(2000), rand.Intn(2000))
					dest.VHDFileChecksum = fmt.Sprintf("%d%d%d", rand.Intn(2000), rand.Intn(2000), rand.Intn(2000))
					expectedDestVHDFileChecksum = dest.VHDFileChecksum
				})

				It("retains dest.VHDFileChecksum's original value", func() {
					Expect(dest.VHDFileChecksum).To(Equal(expectedDestVHDFileChecksum))
				})
			})

			Context("when src does not specify a VHDFileChecksum and dest does not specify a VHDFileChecksum", func() {
				It("should do nothing", func() {
					Expect(dest.VHDFileChecksum).To(BeEmpty())
				})
			})

			Context("when dest does specify a VHDFileChecksum and src does not specify a VHDFileChecksum", func() {
				var expectedDestVHDFileChecksum string

				BeforeEach(func() {
					dest.VHDFileChecksum = fmt.Sprintf("orange%d.vmdk", rand.Intn(2000))
					expectedDestVHDFileChecksum = dest.VHDFileChecksum
				})
				It("should do nothing", func() {
					Expect(dest.VHDFileChecksum).To(Equal(expectedDestVHDFileChecksum))
				})
			})
		})

		Context("Multiple fields", func() {
			Context("when some fields are set in src and another, somewhat overlapping, set of fields is set in dest", func() {
				BeforeEach(func() {
					dest.OutputDir = "needful"
					dest.Version = "do"
					dest.VHDFile = "qwerty"
					dest.VMDKFile = "not"
					src.PatchFile = "matter"
					src.OSVersion = "the"
					src.OutputDir = "bear"
					src.VHDFile = "does"
					src.VHDFileChecksum = "123645125867"
					src.PatchFileChecksum = "123645125867"
				})

				It("copies into dest only those fields which are empty in dest", func() {
					expected := StembuildOptions{
						PatchFile:         "matter",
						OSVersion:         "the",
						OutputDir:         "needful",
						Version:           "do",
						VHDFile:           "qwerty",
						VMDKFile:          "not",
						VHDFileChecksum:   "123645125867",
						PatchFileChecksum: "123645125867",
					}
					Expect(dest).To(Equal(expected))
				})
			})
		})
	})

	Context("LoadOptionsFromManifest", func() {
		var (
			testFileName string
			args         StembuildOptions
			executeErr   error
		)

		JustBeforeEach(func() {
			executeErr = LoadOptionsFromManifest(testFileName, &args)
		})

		BeforeEach(func() {
			args = StembuildOptions{}
		})

		Context("when the file does not exist", func() {
			BeforeEach(func() {
				testFileName = "imagination"
			})

			It("throws an appropriate error", func() {
				Expect(executeErr).To(HaveOccurred())
			})
		})

		Context("when the file exists", func() {
			Context("when the file cannot be read", func() {
				BeforeEach(func() {
					var err error
					testFileName, err = ioutil.TempDir("", "")
					Expect(err).NotTo(HaveOccurred())
				})

				AfterEach(func() {
					Expect(os.RemoveAll(testFileName)).To(Succeed())
				})

				It("throws an appropriate error", func() {
					Expect(executeErr).To(HaveOccurred())
				})
			})

			Context("when the file can be read", func() {
				Context("when the file is not proper YAML", func() {
					BeforeEach(func() {
						testFileName = filepath.Join("..", "testdata", "invalid-yml.yml")
					})

					It("throws a parsing error", func() {
						Expect(executeErr).To(HaveOccurred())
					})
				})

				Context("when the file is proper YAML", func() {
					BeforeEach(func() {
						testFileName = filepath.Join("..", "testdata", "valid-apply-patch.yml")
					})

					It("copies into the arguments the values from the manifest", func() {
						Expect(executeErr).NotTo(HaveOccurred())
						expected := StembuildOptions{
							PatchFile:         "testdata/diff.patch",
							Version:           "1200.0",
							VHDFile:           "testdata/original.vhd",
							VHDFileChecksum:   "246616016f66ad2275364be1a2f625758a963a497ea4d1a1103a1a840c3ef274",
							PatchFileChecksum: "d802a5077d747a4ce36e7318b262714dd01be78b645acab30fc01a2131184b09",
						}
						Expect(args).To(Equal(expected))
					})
				})
			})
		})
	})
})
