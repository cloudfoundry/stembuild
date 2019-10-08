package construct_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/cloudfoundry-incubator/stembuild/remotemanager"
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
			err := CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

			Eventually(session, 20*time.Minute).Should(Exit(0))
			Eventually(session.Out).Should(Say(`mock stemcell automation script executed`))
		})

		It("extracts the WinRM BOSH powershell script and executes it successfully on the guest VM", func() {
			err := CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

			Eventually(session, 20*time.Minute).Should(Exit(0))
			Eventually(session.Out).Should(Say(`Attempting to enable WinRM on the guest vm...WinRm enabled on the guest VM`))

		})

		It("handles special characters", func() {
			isAlphaNumeric, err := regexp.Compile("[a-zA-Z0-9]+")
			Expect(err).ToNot(HaveOccurred())

			if isAlphaNumeric.MatchString(conf.VCenterUsername) && isAlphaNumeric.MatchString(conf.VCenterPassword) {
				Skip("vCenter username or password must contain special characters")
			}
			err = CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

			Eventually(session, 20).Should(Exit(0))
			Eventually(session.Out).Should(Say(`mock stemcell automation script executed`))
		})

		AfterEach(func() {
			rm := remotemanager.NewWinRM(conf.TargetIP, conf.VMUsername, conf.VMPassword)
			err := rm.ExecuteCommand("powershell.exe Remove-Item c:\\provision -recurse")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	It("fails with an appropriate error when LGPO is missing", func() {
		session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

		Eventually(session, 20).Should(Exit(1))
		Eventually(session.Err).Should(Say(`Could not find LGPO.zip in the current directory`))
	})
	It("fails with the appropriate error when the Stembuild Version does not match the Guest OS Version", func() {
		err := CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
		Expect(err).ToNot(HaveOccurred())

		wrongVersionStembuildExecutable, err := helpers.BuildStembuild("dev.1")
		Expect(err).ToNot(HaveOccurred())

		session := helpers.Stembuild(wrongVersionStembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

		Eventually(session, 20).Should(Exit(1))
		Eventually(session.Err).Should(Say("OS version of stembuild and guest OS VM do not match"))
	})

	AfterEach(func() {
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
