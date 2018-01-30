package integration

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	"github.com/pivotal-cf-experimental/stembuild/helpers"
)

var _ = Describe("Apply Patch", func() {
	Context("when valid manifest file", func() {
		AfterEach(func() {
			Expect(os.Remove("bosh-stemcell-2012R2-vsphere-esxi-windows2012R2-go_agent.tgz")).To(Succeed())
		})

		It("creates a stemcell", func() {
			session := helpers.Stembuild("apply-patch", filepath.Join("testdata", "valid-apply-patch.yml"))
			Eventually(session, 5).Should(Exit(0))
			Eventually(session).Should(Say(`created stemcell: .*\.tgz`))
		})
	})
})
