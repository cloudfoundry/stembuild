package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	"github.com/pivotal-cf-experimental/stembuild/helpers"
)

var _ = Describe("Convert VMDK", func() {
	Context("when valid vmdk file", func() {
		var stemcellFilename string
		var inputVmdk string

		Context("stembuild when executed", func() {
			var osVersion string
			var version string
			BeforeEach(func() {
				osVersion = "2012R2"
				version = "1200.0"
				stemcellFilename = fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz", version, osVersion)
				inputVmdk = filepath.Join("..", "testdata", "expected.vmdk")
			})

			AfterEach(func() {
				Expect(os.Remove(stemcellFilename)).To(Succeed())
			})

			It("creates a valid stemcell", func() {
				session := helpers.Stembuild("-vmdk", inputVmdk, "-v", version)
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
	})
})
