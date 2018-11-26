package integration

import (
	"fmt"
	"github.com/cloudfoundry-incubator/stembuild/test/helpers"
	"github.com/cloudfoundry-incubator/stembuild/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Version flag", func() {
	Context("when version provided", func() {
		expectedVersion := fmt.Sprintf("stembuild version %s, Windows Stemcell Building Tool", version.Version)

		It("prints version information", func() {
			session := helpers.Stembuild("--version")

			Eventually(session, 20).Should(Exit(0))
			Eventually(session).Should(Say(expectedVersion))
		})

		It("with command, prints version information and does not run command", func() {
			session := helpers.Stembuild("--version", "package")

			Eventually(session, 20).Should(Exit(0))
			Eventually(session).Should(Say(expectedVersion))
		})
	})
})
