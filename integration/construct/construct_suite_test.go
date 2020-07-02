package construct_test

import (
	"fmt"
	"github.com/cloudfoundry-incubator/stembuild/remotemanager"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cloudfoundry-incubator/stembuild/test/helpers"

	"github.com/concourse/pool-resource/out"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vmware/govmomi/govc/cli"
	_ "github.com/vmware/govmomi/govc/device"
	_ "github.com/vmware/govmomi/govc/importx"
	_ "github.com/vmware/govmomi/govc/vm"
	_ "github.com/vmware/govmomi/govc/vm/guest"
	_ "github.com/vmware/govmomi/govc/vm/snapshot"

	"syscall"
)

func TestConstruct(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Construct Suite")
}

const (
	VMNameVariable                    = "VM_NAME"
	VMUsernameVariable                = "VM_USERNAME"
	VMPasswordVariable                = "VM_PASSWORD"
	TargetVmIPVariable                = "TARGET_VM_IP"
	SkipCleanupVariable               = "SKIP_CLEANUP"
	vcenterFolderVariable             = "VM_FOLDER"
	vcenterAdminCredentialUrlVariable = "VCENTER_ADMIN_CREDENTIAL_URL"
	vcenterBaseURLVariable            = "VCENTER_BASE_URL"
	vcenterStembuildUsernameVariable  = "VCENTER_USERNAME"
	vcenterStembuildPasswordVariable  = "VCENTER_PASSWORD"
	StembuildVersionVariable          = "STEMBUILD_VERSION"
	BoshPsmodulesRepoVariable         = "BOSH_PSMODULES_REPO"
	VmSnapshotName                    = "integration-test-snapshot"
	LoggedInVmIpVariable              = "LOGOUT_INTEGRATION_TEST_VM_IP"
	LoggedInVmIpathVariable           = "LOGOUT_INTEGRATION_TEST_VM_INVENTORY_PATH"
	LoggedInVmSnapshotName            = "logged-in-winrm-enabled"
	powershell                        = "C:\\Windows\\System32\\WindowsPowerShell\\V1.0\\powershell.exe"
)

var (
	conf                      config
	tmpDir                    string
	lockParentDir             string
	lockPool                  out.LockPool
	lockDir                   string
	stembuildExecutable       string
	vcenterAdminCredentialUrl string
)

type config struct {
	TargetIP           string
	NetworkGateway     string
	SubnetMask         string
	VMUsername         string
	VMPassword         string
	VMName             string
	VMNetwork          string
	VCenterURL         string
	VCenterUsername    string
	VCenterPassword    string
	VMInventoryPath    string
	LoggedInVMIP       string
	LoggedInVMIpath    string
	LoggedInVMSnapshot string
}

var _ = SynchronizedBeforeSuite(func() []byte {
	var err error

	boshPsmodulesRepo := envMustExist(BoshPsmodulesRepoVariable)
	stembuildVersion := envMustExist(StembuildVersionVariable)
	stembuildExecutable, err = helpers.BuildStembuild(stembuildVersion)
	Expect(err).NotTo(HaveOccurred())

	vmUsername := envMustExist(VMUsernameVariable)
	vmPassword := envMustExist(VMPasswordVariable)
	targetVMIP := envMustExist(TargetVmIPVariable)
	vmName := envMustExist(VMNameVariable)

	loggedInVmIp := envMustExist(LoggedInVmIpVariable)
	loggedInVmInventoryPath := envMustExist(LoggedInVmIpathVariable)
	loggedInVmSnapshot := LoggedInVmSnapshotName
	vCenterUrl := envMustExist(vcenterBaseURLVariable)
	vcenterFolder := envMustExist(vcenterFolderVariable)
	vmInventoryPath := strings.Join([]string{vcenterFolder, vmName}, "/")
	vcenterAdminCredentialUrl = envMustExist(vcenterAdminCredentialUrlVariable)

	vCenterStembuildUser := envMustExist(vcenterStembuildUsernameVariable)
	vCenterStembuildPassword := envMustExist(vcenterStembuildPasswordVariable)

	wd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())
	tmpDir, err = ioutil.TempDir(wd, "construct-integration")
	Expect(err).NotTo(HaveOccurred())

	err = os.MkdirAll(tmpDir, 0755)
	Expect(err).NotTo(HaveOccurred())

	conf = config{
		TargetIP:           targetVMIP,
		VMUsername:         vmUsername,
		VMPassword:         vmPassword,
		VCenterURL:         vCenterUrl,
		VCenterUsername:    vCenterStembuildUser,
		VCenterPassword:    vCenterStembuildPassword,
		LoggedInVMIP:       loggedInVmIp,
		LoggedInVMIpath:    loggedInVmInventoryPath,
		LoggedInVMSnapshot: loggedInVmSnapshot,
		VMName:             vmName,
		VMInventoryPath:    vmInventoryPath,
	}

	enableWinRM(boshPsmodulesRepo)
	powerOnVM()
	createVMSnapshot(VmSnapshotName)

	return nil
}, func(_ []byte) {
})

var _ = BeforeEach(func() {
	revertSnapshot(conf.VMInventoryPath, VmSnapshotName)
	waitForVmToBeReady(conf.TargetIP, conf.VMUsername, conf.VMPassword)
})

var _ = SynchronizedAfterSuite(func() {
	skipCleanup := strings.ToUpper(os.Getenv(SkipCleanupVariable))

	if skipCleanup != "TRUE" {
		deleteCommand := []string{
			"vm.destroy",
			fmt.Sprintf("-vm.ipath=%s", conf.VMInventoryPath),
			fmt.Sprintf("-u=%s", vcenterAdminCredentialUrl),
		}
		Eventually(func() int {
			return cli.Run(deleteCommand)
		}, 3*time.Minute, 10*time.Second).Should(BeZero())
		fmt.Println("VM destroyed")
		if lockDir != "" {
			_, _, err := lockPool.ReleaseLock(lockDir)
			Expect(err).NotTo(HaveOccurred())

			childItems, err := ioutil.ReadDir(lockParentDir)
			Expect(err).NotTo(HaveOccurred())

			for _, item := range childItems {
				if item.IsDir() && strings.HasPrefix(filepath.Base(item.Name()), "pool-resource") {
					fmt.Printf("Cleaning up temporary pool resource %s\n", item.Name())
					_ = os.RemoveAll(item.Name())
				}
			}
		}
	}

	_ = os.RemoveAll(tmpDir)
}, func() {
})

func revertSnapshot(vmIpath string, snapshotName string) {
	snapshotCommand := []string{
		"snapshot.revert",
		fmt.Sprintf("-vm.ipath=%s", vmIpath),
		fmt.Sprintf("-u=%s", vcenterAdminCredentialUrl),
		snapshotName,
	}
	fmt.Printf("Reverting VM Snapshot: %s\n", snapshotName)
	exitCode := runIgnoringOutput(snapshotCommand)
	if exitCode != 0 {
		fmt.Print("There was an error reverting the snapshot.")
	} else {
		fmt.Println("Revert started.")
	}
	time.Sleep(30 * time.Second)
}

func waitForVmToBeReady(vmIp string, vmUsername string, vmPassword string) {
	fmt.Print("Waiting for reverting snapshot to finish...")
	clientFactory := remotemanager.NewWinRmClientFactory(vmIp, vmUsername, vmPassword)
	rm := remotemanager.NewWinRM(vmIp, vmUsername, vmPassword, clientFactory)
	Expect(rm).ToNot(BeNil())

	vmReady := false
	for !vmReady {
		fmt.Print(".")
		time.Sleep(5 * time.Second)
		_, err := rm.ExecuteCommand(`powershell.exe "ls c:\windows 1>$null"`)
		vmReady = err == nil
	}
	fmt.Print("done.\n")
}

func envMustExist(variableName string) string {
	result := os.Getenv(variableName)
	if result == "" {
		Fail(fmt.Sprintf("%s must be set", variableName))
	}

	return result
}

func enableWinRM(repoPath string) {
	fmt.Println("Enabling WinRM on the base image before integration tests...")
	uploadCommand := []string{
		"guest.upload",
		fmt.Sprintf("-vm.ipath=%s", conf.VMInventoryPath),
		fmt.Sprintf("-u=%s", vcenterAdminCredentialUrl),
		fmt.Sprintf("-l=%s:%s", conf.VMUsername, conf.VMPassword),
		filepath.Join(repoPath, "modules", "BOSH.WinRM", "BOSH.WinRM.psm1"),
		"C:\\Windows\\Temp\\BOSH.WinRM.psm1",
	}

	exitCode := runIgnoringOutput(uploadCommand)
	if exitCode != 0 {
		fmt.Print("There was an error uploading WinRM psmodule.")
	}

	enableCommand := []string{
		"guest.start",
		fmt.Sprintf("-vm.ipath=%s", conf.VMInventoryPath),
		fmt.Sprintf("-u=%s", vcenterAdminCredentialUrl),
		fmt.Sprintf("-l=%s:%s", conf.VMUsername, conf.VMPassword),
		powershell,
		`-command`,
		`&{Import-Module C:\Windows\Temp\BOSH.WinRM.psm1; Enable-WinRM}`,
	}
	exitCode = runIgnoringOutput(enableCommand)
	if exitCode != 0 {
		fmt.Print("There was an error enabling WinRM.")
	} else {
		fmt.Println("WinRM enabled.")
	}
}

func createVMSnapshot(snapshotName string) {
	snapshotCommand := []string{
		"snapshot.create",
		fmt.Sprintf("-vm.ipath=%s", conf.VMInventoryPath),
		fmt.Sprintf("-u=%s", vcenterAdminCredentialUrl),
		snapshotName,
	}
	fmt.Printf("Creating VM Snapshot: %s on VM: %s\n", snapshotName, conf.VMInventoryPath)
	// is blocking
	exitCode := runIgnoringOutput(snapshotCommand)
	if exitCode != 0 {
		fmt.Print("There was an error creating the snapshot.")
	} else {
		fmt.Println("Snapshot command started.")
	}
	fmt.Print("Waiting for snapshot to finish...")
	time.Sleep(30 * time.Second)
	fmt.Print("done.\n")
}

func powerOnVM() {
	powerOnCommand := []string{
		"vm.power",
		fmt.Sprintf("-vm.ipath=%s", conf.VMInventoryPath),
		fmt.Sprintf("-u=%s", vcenterAdminCredentialUrl),
		fmt.Sprintf("-on"),
	}
	runIgnoringOutput(powerOnCommand)
}

func runIgnoringOutput(args []string) int {
	oldStderr := os.Stderr
	oldStdout := os.Stdout

	_, w, _ := os.Pipe()

	defer w.Close()

	os.Stderr = w
	os.Stdout = w

	os.Stderr = os.NewFile(uintptr(syscall.Stderr), "/dev/null")

	exitCode := cli.Run(args)

	os.Stderr = oldStderr
	os.Stdout = oldStdout

	return exitCode
}
