package integration_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	"strings"

	"time"

	"github.com/cloudfoundry-incubator/stembuild/test/helpers"
)

var _ = Describe("Convert VMDK", func() {

	Context("when valid vmdk file", func() {
		var stemcellFilename string
		var inputVmdk string

		Context("stembuild when executed with invalid", func() {
			var osVersion string
			var version string

			Context("OS value", func() {
				It("of 1709 returns an error", func() {
					osVersion = "1709"
					version = "1709.0"
					expectedOSVersionInNameANdManifest := "2016"

					stemcellFilename = fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz", version, expectedOSVersionInNameANdManifest)
					inputVmdk = filepath.Join("..", "test", "data", "expected.vmdk")

					session := helpers.Stembuild(stembuildExecutable, "package", "--vmdk", inputVmdk, "--stemcell-version", version, "--os", osVersion)
					Eventually(session, 20).Should(Exit(1))
					Eventually(session.Err).Should(Say(`OS version must be either 2012R2, 2016, 1803, or 2019 have:`))
				})

			})
		})

		Context("stembuild when executed", func() {
			var osVersion string
			var version string

			AfterEach(func() {
				Expect(os.Remove(stemcellFilename)).To(Succeed())
			})

			It("creates a valid 2012R2 stemcell", func() {
				osVersion = "2012R2"
				version = "1200.0"
				stemcellFilename = fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz", version, osVersion)
				inputVmdk = filepath.Join("..", "test", "data", "expected.vmdk")

				session := expectStembuildToSucceed("package", "--vmdk", inputVmdk, "--stemcell-version", version, "--os", osVersion, "--outputDir", ".")
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

			It("creates a valid 1803 stemcell", func() {
				osVersion = "1803"
				version = "1803.0"
				stemcellFilename = fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz", version, osVersion)
				inputVmdk = filepath.Join("..", "test", "data", "expected.vmdk")

				session := expectStembuildToSucceed("package", "--vmdk", inputVmdk, "--stemcell-version", version, "--os", osVersion, "--outputDir", ".")
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

			It("creates a valid 2019 stemcell", func() {
				osVersion = "2019"
				version = "2019.0"
				stemcellFilename = fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz", version, osVersion)
				inputVmdk = filepath.Join("..", "test", "data", "expected.vmdk")

				session := expectStembuildToSucceed("package", "--vmdk", inputVmdk, "--stemcell-version", version, "--os", osVersion, "--outputDir", ".")
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

			It("creates a valid 2016 stemcell", func() {
				osVersion = "2016"
				version = "2016.0"
				stemcellFilename = fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz", version, osVersion)
				inputVmdk = filepath.Join("..", "test", "data", "expected.vmdk")

				session := expectStembuildToSucceed("package", "--vmdk", inputVmdk, "--stemcell-version", version, "--os", osVersion, "--outputDir", ".")

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
	})
})

func expectStembuildToSucceed(arguments ...string) *Session {
	session := helpers.Stembuild(stembuildExecutable, arguments...)
	Eventually(session, 20*time.Second).Should(Exit(0),
		fmt.Sprintf(
			"Expected %s %s to exit with code 0, exited with code %d\nout: %s\nerr: %s",
			stembuildExecutable,
			strings.Join(arguments, " "),
			session.ExitCode(),
			string(session.Out.Contents()),
			string(session.Err.Contents()),
		))

	return session
}
