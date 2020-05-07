package construct_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/cloudfoundry-incubator/stembuild/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

const (
	vCenterUsername = "USER"
	vCenterPassword = "PASS"
)

var _ = Describe("stembuild construct", func() {
	var workingDir string

	BeforeEach(func() {
		var err error
		workingDir, err = os.Getwd()
		Expect(err).ToNot(HaveOccurred())

	})

	const constructOutputTimeout = 60 * time.Second
	Context("run successfully", func() {

		FIt("successfully exits when vm becomes powered off", func() {
			err := CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

			shutdownTimeout := 3 * time.Minute
			Eventually(session, shutdownTimeout).Should(Exit(0))
		})

		It("transfers LGPO and StemcellAutomation archives, unarchive them and execute automation script", func() {
			err := CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

			Eventually(session.Out, constructOutputTimeout).Should(Say(`mock stemcell automation script executed`))
		})

		It("executes post-reboot automation script", func() {
			err := CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

			Eventually(session.Out, constructOutputTimeout*5).Should(Say(`mock stemcell automation post-reboot script executed`))
		})

		It("extracts the WinRM BOSH powershell script and executes it successfully on the guest VM", func() {
			err := CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

			Eventually(session.Out, constructOutputTimeout).Should(Say(`Attempting to enable WinRM on the guest vm...WinRm enabled on the guest VM`))

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

			Eventually(session, constructOutputTimeout).Should(Exit(0))
			Eventually(session.Out).Should(Say(`mock stemcell automation script executed`))
		})

		It("successfully runs even when a user has logged in", func() {
			loggedOnVmIp := "10.74.34.88"
			loggedOnVmInventoryPath := "/canada-dc/vm/calgary/construct-integration-test88"
			loggedOnVmSnapshot := "logged-on"

			// revert snapshot
			snapshotCommand := []string{
				"snapshot.revert",
				fmt.Sprintf("-vm.ipath=%s", loggedOnVmInventoryPath),
				fmt.Sprintf("-u=%s", vcenterAdminCredentialUrl),
				loggedOnVmSnapshot,
			}

			fmt.Printf("Reverting VM Snapshot: %s", loggedOnVmSnapshot)
			runIgnoringOutput(snapshotCommand)

			//run normal stembuild construct command, like we do in prev. test
			err := CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "construct",
				"-vm-ip", loggedOnVmIp,
				"-vm-username", conf.VMUsername,
				"-vm-password", conf.VMPassword,
				"-vcenter-url", conf.VCenterURL,
				"-vcenter-username", conf.VCenterUsername,
				"-vcenter-password", conf.VCenterPassword,
				"-vm-inventory-path", loggedOnVmInventoryPath)

			shutdownTimeout := 3 * time.Minute
			Eventually(session, shutdownTimeout).Should(Exit(0))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("fails with an appropriate error when LGPO is missing", func() {
		session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

		Eventually(session, constructOutputTimeout).Should(Exit(1))
		Eventually(session.Err).Should(Say(`Could not find LGPO.zip in the current directory`))
	})

	It("does not exit when the target VM has not powered off", func() {
		err := CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
		Expect(err).ToNot(HaveOccurred())

		fakeStemcellAutomationShutdownDelay := 45 * time.Second

		session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

		Consistently(session, fakeStemcellAutomationShutdownDelay-5*time.Second).Should(Not(Exit()))
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
