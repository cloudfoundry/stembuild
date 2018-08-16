package integration

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"

	"github.com/pivotal-cf-experimental/stembuild/helpers"
	"github.com/pivotal-cf-experimental/stembuild/stembuildoptions"
)

const validManifestTemplate = helpers.ManifestTemplate
const invalidManifestTemplate = `---
version: "2012R2"
dhv_flie_nmae "some-vhd-file"
ptach_flie: "some-patch-file"
`

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
		var (
			stemcellFilename string
			osVersion        string
		)

		BeforeEach(func() {
			manifestStruct.Version = "1200.0"
			manifestStruct.VHDFile = "testdata/original.vhd"
			manifestStruct.PatchFile = "testdata/diff.patch"
			manifestStruct.VHDFileChecksum = "246616016f66ad2275364be1a2f625758a963a497ea4d1a1103a1a840c3ef274"
			manifestStruct.PatchFileChecksum = "d802a5077d747a4ce36e7318b262714dd01be78b645acab30fc01a2131184b09"
			manifestText = validManifestTemplate
			manifestStruct.OSVersion = "2012R2"
			osVersion = "2012R2"
			stemcellFilename = fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz", manifestStruct.Version, osVersion)
		})

		Context("stembuild when executed with a patchfile on disk", func() {
			BeforeEach(func() {
				manifestStruct.OSVersion = "2016"
				osVersion = "2016"
				stemcellFilename = fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz", manifestStruct.Version, osVersion)
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

		Context("when vhd file checksum does not match the checksum in the manifest file", func() {
			BeforeEach(func() {
				manifestStruct.VHDFileChecksum = "incorrect checksum value"
			})
			It("fails with the expected error message", func() {
				session := helpers.Stembuild("apply-patch", manifestFilename)
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("the specified base VHD is different from the VHD expected by the diff bundle"))
			})
		})

		Context("when patch file checksum does not match the checksum in the manifest file", func() {
			BeforeEach(func() {
				manifestStruct.PatchFileChecksum = "incorrect checksum value"
			})
			It("fails with the expected error message", func() {
				session := helpers.Stembuild("apply-patch", manifestFilename)
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("the specified patch file is different from the patch file expected by the diff bundle"))
			})
		})

		Context("when OS version is invalid", func() {
			BeforeEach(func() {
				manifestStruct.OSVersion = ""
			})

			It("displays an error", func() {
				session := helpers.Stembuild("apply-patch", manifestFilename)
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("OS version must be either 2012R2, 2016 or 1803"))
			})
		})

		Context("when OS version is 1803", func() {
			BeforeEach(func() {
				manifestStruct.OSVersion = "1803"
				osVersion = "1803"
				manifestStruct.Version = "1803.0"
				stemcellFilename = fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz", manifestStruct.Version, osVersion)
			})

			AfterEach(func() {
				Expect(os.Remove(stemcellFilename)).To(Succeed())
			})

			It("creates a valid stemcell", func() {
				session := helpers.Stembuild("apply-patch", manifestFilename)
				Eventually(session, 5).Should(Exit(0))

				stemcellDir, err := helpers.ExtractGzipArchive(stemcellFilename)
				Expect(err).NotTo(HaveOccurred())
				manifestFilepath := filepath.Join(stemcellDir, "stemcell.MF")
				manifest, err := helpers.ReadFile(manifestFilepath)
				Expect(err).NotTo(HaveOccurred())
				expectedName := fmt.Sprintf("name: bosh-vsphere-esxi-windows%s-go_agent", osVersion)
				Expect(manifest).To(ContainSubstring(expectedName))
			})
		})

		Context("when no output directory is specified on the command line", func() {
			Context("current working directory has no stemcell tgz in it", func() {
				AfterEach(func() {
					Expect(os.Remove(stemcellFilename)).To(Succeed())
				})

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

				AfterEach(func() {
					Expect(os.Remove(stemcellFilename)).To(Succeed())
				})

				It("displays an error", func() {
					session := helpers.Stembuild("apply-patch", manifestFilename)
					Eventually(session).Should(Exit(1))
					Eventually(session.Err).Should(Say("file may already exist"))
					Eventually(session.Err).Should(Say(`\n\nfor usage: stembuild -h`))
				})
			})

			Context("when stembuild is executed with a url pointing at a patchfile", func() {
				var patchServer *Server

				BeforeEach(func() {
					var patchURL string
					patchServer, patchURL = helpers.StartFileServer("testdata/diff.patch")
					manifestStruct.PatchFile = patchURL
				})

				AfterEach(func() {
					patchServer.Close()
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

			Context("the patchfile url is invalid", func() {
				var (
					patchServer *Server
					patchPath   string
				)

				BeforeEach(func() {
					patchServer, patchPath = helpers.StartInvalidFileServer(http.StatusNotFound)
					manifestStruct.PatchFile = patchPath
				})

				AfterEach(func() {
					patchServer.Close()
				})

				It("fails and returns an error", func() {
					session := helpers.Stembuild("apply-patch", manifestFilename)
					Eventually(session).Should(Exit(1))
					Eventually(session.Err).Should(Say(`Error: Could not create stemcell from %s`, patchPath))
					Eventually(session.Err).Should(Say("Unexpected response code: %d", http.StatusNotFound))
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
