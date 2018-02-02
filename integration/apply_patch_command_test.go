package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	"github.com/pivotal-cf-experimental/stembuild/helpers"
	"github.com/pivotal-cf-experimental/stembuild/stembuildoptions"
)

var _ = Describe("Apply Patch", func() {
	var manifestStruct stembuildoptions.StembuildOptions
	var manifestText string
	var manifestFilename string
	BeforeEach(func() {
		manifestStruct = stembuildoptions.StembuildOptions{}
	})
	JustBeforeEach(func() {
		manifestFile, err := ioutil.TempFile("", "")
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			Expect(manifestFile.Close()).To(Succeed())
		}()

		contents, err := helpers.StringFromManifest(manifestText, manifestStruct)
		Expect(err).NotTo(HaveOccurred())
		_, err = manifestFile.Write([]byte(contents))
		Expect(err).NotTo(HaveOccurred())

		manifestFilename = manifestFile.Name()
	})

	Context("when valid manifest file", func() {
		var stemcellFilename string
		const validManifestTemplate = helpers.ManifestTemplate

		BeforeEach(func() {
			manifestStruct.Version = "1200.0"
			manifestStruct.VHDFile = "testdata/original.vhd"
			manifestStruct.PatchFile = "testdata/diff.patch"
			manifestText = validManifestTemplate
		})

		Context("stembuild when executed", func() {
			var osVersion string
			BeforeEach(func() {
				osVersion = "2012R2"
				stemcellFilename = fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz", manifestStruct.Version, osVersion)
				manifestStruct.VHDFile = "testdata/original.vhd"
				manifestStruct.PatchFile = "testdata/diff.patch"
			})

			AfterEach(func() {
				Expect(os.Remove(stemcellFilename)).To(Succeed())
			})

			It("creates a valid stemcell", func() {
				session := helpers.Stembuild("apply-patch", manifestFilename)
				Eventually(session, 5).Should(Exit(0))
				Eventually(session).Should(Say(`created stemcell: .*\.tgz`))
				Expect(stemcellFilename).To(BeAnExistingFile())

				stemcellDir, err := helpers.ExtractGzipArchive(stemcellFilename)
				Expect(err).NotTo(HaveOccurred())

				manifestFilepath := filepath.Join(stemcellDir, "stemcell.MF")
				manifest, err := helpers.ReadFile(manifestFilepath)
				Expect(err).NotTo(HaveOccurred())

				expectedOs := fmt.Sprintf("operating_system: windows%s", osVersion)
				Expect(manifest).To(ContainSubstring(expectedOs))

				expectedName := fmt.Sprintf("name: bosh-vsphere-esxi-windows%s-go_agent", osVersion)
				Expect(manifest).To(ContainSubstring(expectedName))

				imageFilepath := filepath.Join(stemcellDir, "image")
				imageDir, err := helpers.ExtractGzipArchive(imageFilepath)
				Expect(err).NotTo(HaveOccurred())

				actualVmdkFilepath := filepath.Join(imageDir, "image-disk1.vmdk")
				_, err = ioutil.ReadFile(actualVmdkFilepath)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when no output directory is specified on the command line", func() {
			BeforeEach(func() {
				osVersion := "2012R2"
				stemcellFilename = fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz", manifestStruct.Version, osVersion)
				manifestStruct.VHDFile = "testdata/original.vhd"
				manifestStruct.PatchFile = "testdata/diff.patch"
			})

			AfterEach(func() {
				Expect(os.Remove(stemcellFilename)).To(Succeed())
			})

			Context("current working directory has no stemcell tgz in it", func() {
				It("creates a stemcell in current working directory", func() {
					session := helpers.Stembuild("apply-patch", manifestFilename)
					Eventually(session, 5).Should(Exit(0))
					Eventually(session).Should(Say(`created stemcell: .*\.tgz`))
				})
			})

			Context("current working directory has stemcell tgz in it", func() {
				BeforeEach(func() {
					stemcellFile, err := os.Create(stemcellFilename)
					Expect(err).NotTo(HaveOccurred())
					stemcellFile.Close()
				})

				It("displays an error", func() {
					session := helpers.Stembuild("apply-patch", manifestFilename)
					Eventually(session).Should(Exit(1))
					Eventually(session.Err).Should(Say("file may already exist"))
					Eventually(session.Err).Should(Say(`\n\nfor usage: stembuild -h`))
				})
			})
		})

		Context("when output directory specified with -o flag", func() {
			var tmpDir string

			BeforeEach(func() {
				var err error
				tmpDir, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				Expect(os.RemoveAll(tmpDir)).To(Succeed())
			})

			Context("directory already exists", func() {
				// TODO: what if the directory already has a stemcell in it?
				It("creates stemcell in output directory", func() {
					session := helpers.Stembuild("-o", tmpDir, "apply-patch", manifestFilename)
					Eventually(session, 5).Should(Exit(0))
					safeDir := strings.Replace(tmpDir, `\`, `\\`, -1)
					Eventually(session).Should(Say(`created stemcell: .*%s.*\.tgz`, safeDir))
				})
			})

			Context("directory does not exist", func() {
				AfterEach(func() {
					Expect(os.RemoveAll("idontexist")).To(Succeed())
				})
				It("creates directory and puts stemcell in it", func() {
					session := helpers.Stembuild("-o", "idontexist", "apply-patch", manifestFilename)
					Eventually(session, 5).Should(Exit(0))
					Eventually(session).Should(Say(`created stemcell: .*idontexist.*\.tgz`))
					Expect(helpers.Exists("idontexist")).To(BeTrue())
				})
			})
		})
	})

	Context("Invalid apply-patch manifest file", func() {
		invalidManifestTemplate := `---
		version: "2012R2"
		dhv_flie_nmae "some-vhd-file"
		ptach_flie: "some-patch-file"
		`
		BeforeEach(func() {
			manifestStruct.Version = "1200.0"
			manifestStruct.VHDFile = "testdata/original.vhd"
			manifestStruct.PatchFile = "testdata/diff.patch"
			manifestText = invalidManifestTemplate
		})
		It("Returns an error", func() {
			session := helpers.Stembuild("apply-patch", manifestFilename)
			Eventually(session).Should(Exit(1))
		})
	})
})
