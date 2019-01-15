package construct_test

import (
	"fmt"
	"github.com/cloudfoundry-incubator/stembuild/remotemanager"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/stembuild/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("stembuild construct", func() {
	var workingDir string

	BeforeEach(func() {
		var err error
		workingDir, err = os.Getwd()
		Expect(err).ToNot(HaveOccurred())

	})

	Context("run successfully", func() {
		It("transfers LGPO and StemcellAutomation archives, unarchive them and execute automation script", func() {
			err := CopyFile(filepath.Join(workingDir, "assets", "StemcellAutomation.zip"), filepath.Join(workingDir, "StemcellAutomation.zip"))
			Expect(err).ToNot(HaveOccurred())

			err = CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "construct", "-winrm-ip", conf.TargetIP, "-stemcell-version", "1709.1", "-winrm-username", conf.VMUsername, "-winrm-password", conf.VMPassword)

			Eventually(session, 20).Should(Exit(0))
			Eventually(session.Out).Should(Say(`mock stemcell automation script executed`))
		})

		AfterEach(func() {
			rm := remotemanager.NewWinRM(conf.TargetIP, conf.VMUsername, conf.VMPassword)
			err := rm.ExecuteCommand("powershell.exe Remove-Item c:\\provision -recurse")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	It("fails with an appropriate error when LGPO and/or StemcellAutomation is missing", func() {
		session := helpers.Stembuild(stembuildExecutable, "construct", "-winrm-ip", conf.TargetIP, "-stemcell-version", "1803.1", "-winrm-username", conf.VMUsername, "-winrm-password", conf.VMPassword)

		Eventually(session, 20).Should(Exit(1))
		Eventually(session.Err).Should(Say(`automation artifact not found in current directory`))
	})

	AfterEach(func() {
		_ = os.Remove(filepath.Join(workingDir, "StemcellAutomation.zip"))
		_ = os.Remove(filepath.Join(workingDir, "LGPO.zip"))
	})
})

func CopyFile(src string, dest string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = ioutil.WriteFile(dest, input, 0644)
	if err != nil {
		fmt.Println("Error creating file")
		fmt.Println(err)
		return err
	}

	return err
}
