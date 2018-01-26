package patch_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/stembuild/patch"
)

var _ = Describe("Apply Patch", func() {
	Context("CopyInto", func() {
		var (
			dest ApplyPatch
			src  ApplyPatch
		)

		BeforeEach(func() {
			src = ApplyPatch{}
			dest = ApplyPatch{}
		})

		JustBeforeEach(func() {
			dest.CopyInto(src)
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

			Context("when src specifies a PatchFile and dest specifies a PatchFile", func() {
				It("should do nothing", func() {
					Expect(dest.PatchFile).To(BeEmpty())
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

			Context("when src specifies an OSVersion and dest specifies an OSVersion", func() {
				It("should do nothing", func() {
					Expect(dest.OSVersion).To(BeEmpty())
				})
			})
		})

		Context("OutputDir", func() {
			Context("when src specifies an OutputDir and dest does not", func() {
				BeforeEach(func() {
					src.OutputDir = fmt.Sprintf("foo%d/%d", rand.Intn(2000), rand.Intn(2000))
				})

				It("copies src.OutputDir into dest.OutputDir", func() {
					Expect(dest.OutputDir).To(Equal(src.OutputDir))
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

			Context("when src specifies an OutputDir and dest specifies an OutputDir", func() {
				It("should do nothing", func() {
					Expect(dest.OutputDir).To(BeEmpty())
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

			Context("when src specifies a version and dest specifies a version", func() {
				It("should do nothing", func() {
					Expect(dest.Version).To(BeEmpty())
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

			Context("when src specifies a VHDFile and dest specifies a VHDFile", func() {
				It("should do nothing", func() {
					Expect(dest.VHDFile).To(BeEmpty())
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

			Context("when src specifies a VMDKFile and dest specifies a VMDKFile", func() {
				It("should do nothing", func() {
					Expect(dest.VMDKFile).To(BeEmpty())
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
				})

				It("copies into dest only those fields which are empty in dest", func() {
					expected := ApplyPatch{
						PatchFile: "matter",
						OSVersion: "the",
						OutputDir: "needful",
						Version:   "do",
						VHDFile:   "qwerty",
						VMDKFile:  "not",
					}
					Expect(dest).To(Equal(expected))
				})
			})
		})
	})

	Context("LoadPatchManifest", func() {
		var (
			testFileName string
			args         ApplyPatch
			executeErr   error
		)

		JustBeforeEach(func() {
			executeErr = LoadPatchManifest(testFileName, &args)
		})

		BeforeEach(func() {
			args = ApplyPatch{}
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
						testFileName = filepath.Join("testdata", "invalid-yml.yml")
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
						expected := ApplyPatch{
							PatchFile: "some-patch-file",
							Version:   "2012R2",
							VHDFile:   "some-vhd-file",
						}
						Expect(args).To(Equal(expected))
					})
				})
			})
		})
	})
})
