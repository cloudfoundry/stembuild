package construct_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/cloudfoundry-incubator/stembuild/test/helpers"

	"github.com/masterzen/winrm"

	"github.com/concourse/pool-resource/out"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vmware/govmomi/govc/cli"
	_ "github.com/vmware/govmomi/govc/importx"
	_ "github.com/vmware/govmomi/govc/vm"
)

func TestConstruct(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Construct Suite")
}

const (
	NetworkGatewayVariable = "NETWORK_GATEWAY"
	SubnetMaskVariable     = "SUBNET_MASK"
	OvaFileVariable        = "OVA_FILE"
	VMNamePrefixVariable   = "VM_NAME_PREFIX"
	VMUsernameVariable     = "VM_USERNAME"
	VMPasswordVariable     = "VM_PASSWORD"
	ExistingVmIPVariable   = "EXISTING_VM_IP"
	UserProvidedIPVariable = "USER_PROVIDED_IP"
	LockPrivateKeyVariable = "LOCK_PRIVATE_KEY"
)

var (
	conf                config
	tmpDir              string
	lockPool            out.LockPool
	lockDir             string
	stembuildExecutable string
	existingVM          bool
)

type config struct {
	TargetIP       string
	NetworkGateway string
	SubnetMask     string
	VMUsername     string
	VMPassword     string
	VMName         string
}

func envMustExist(variableName string) string {
	result := os.Getenv(variableName)
	if result == "" {
		Fail(fmt.Sprintf("%s must be set", variableName))
	}

	return result
}

func claimAvailableIP() string {
	lockPrivateKey := envMustExist(LockPrivateKeyVariable)
	keyFile, err := ioutil.TempFile(os.TempDir(), "keyfile")
	Expect(err).NotTo(HaveOccurred())
	_, _ = keyFile.Write([]byte(lockPrivateKey))
	_ = keyFile.Chmod(0600)

	err = exec.Command("ssh-add", keyFile.Name()).Run()
	Expect(err).NotTo(HaveOccurred())

	poolSource := out.Source{
		URI:        "git@github.com:pivotal-cf-experimental/Bosh-Windows-Locks.git",
		Branch:     "master",
		Pool:       "vcenter-ips",
		RetryDelay: 5 * time.Second,
	}

	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)

	lockPool = out.NewLockPool(poolSource, writer)

	ip, _, err := lockPool.AcquireLock()
	Expect(err).NotTo(HaveOccurred())
	Expect(ip).NotTo(Equal(""))

	lockDir, err = ioutil.TempDir("", "acquired-lock")
	Expect(err).NotTo(HaveOccurred())
	err = ioutil.WriteFile(filepath.Join(lockDir, "name"), []byte(ip), os.ModePerm)
	Expect(err).NotTo(HaveOccurred())

	return ip
}

var _ = SynchronizedBeforeSuite(func() []byte {
	existingVM = false
	var err error
	var targetIP string
	stembuildExecutable, err = helpers.BuildStembuild()
	Expect(err).NotTo(HaveOccurred())

	vmUsername := envMustExist(VMUsernameVariable)
	vmPassword := envMustExist(VMPasswordVariable)
	existingVMIP := os.Getenv(ExistingVmIPVariable)
	userProvidedIP := os.Getenv(UserProvidedIPVariable)

	conf = config{
		TargetIP:   existingVMIP,
		VMUsername: vmUsername,
		VMPassword: vmPassword,
	}

	if userProvidedIP != "" && existingVMIP == "" {
		fmt.Printf("Creating VM with IP: %s\n", targetIP)
		targetIP = userProvidedIP
		createVMWithIP(targetIP)
	}
	if existingVMIP != "" {
		existingVM = true
		fmt.Printf("Using existing VM with IP: %s\n", existingVMIP)
		targetIP = existingVMIP
	}
	if targetIP == "" {
		fmt.Println("Finding available IP...")
		targetIP = claimAvailableIP()
		createVMWithIP(targetIP)
	}

	fmt.Println("Attempting to connect to VM")
	endpoint := winrm.NewEndpoint(targetIP, 5985, false, true, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, vmUsername, vmPassword)
	Expect(err).NotTo(HaveOccurred())

	var shell *winrm.Shell
	Eventually(func() error {
		shell, err = client.CreateShell()
		return err
	}, 3*time.Minute).Should(BeNil())
	_ = shell.Close()
	fmt.Println("Successfully connected to VM")

	return nil
}, func(_ []byte) {
})

func createVMWithIP(targetIP string) {
	ovaFile := envMustExist(OvaFileVariable)

	vmNamePrefix := envMustExist(VMNamePrefixVariable)
	conf.NetworkGateway = envMustExist(NetworkGatewayVariable)
	conf.SubnetMask = envMustExist(SubnetMaskVariable)

	conf.TargetIP = targetIP
	fmt.Printf("Target ip is %s\n", targetIP)

	vmNameSuffix := strings.Split(targetIP, ".")[3]
	vmName := fmt.Sprintf("%s%s", vmNamePrefix, vmNameSuffix)
	conf.VMName = vmName

	templateFile, err := filepath.Abs("assets/ova_options.json.template")
	Expect(err).NotTo(HaveOccurred())
	tmpl, err := template.New("ova_options.json.template").ParseFiles(templateFile)

	tmpDir, err = ioutil.TempDir("", "construct-test")
	Expect(err).NotTo(HaveOccurred())

	tmpFile, err := ioutil.TempFile(tmpDir, "ova_options*.json")
	Expect(err).NotTo(HaveOccurred())

	err = tmpl.Execute(tmpFile, conf)
	Expect(err).NotTo(HaveOccurred())

	opts := []string{
		"import.ova",
		fmt.Sprintf("--options=%s", tmpFile.Name()),
		fmt.Sprintf("--name=%s", vmName),
		"--folder=/canada-dc/vm/winnipeg",
		ovaFile,
	}

	fmt.Printf("Opts are %s", opts)

	exitCode := cli.Run(opts)
	Expect(exitCode).To(BeZero())

}

var _ = SynchronizedAfterSuite(func() {
	_ = os.RemoveAll(tmpDir)

	if !existingVM {
		deleteCommand := []string{"vm.destroy", fmt.Sprintf("-vm.ip=%s", conf.TargetIP)}
		Eventually(func() int {
			return cli.Run(deleteCommand)
		}, 3*time.Minute).Should(BeZero())
		fmt.Println("VM destroyed")
		if lockDir != "" {
			_, _, err := lockPool.ReleaseLock(lockDir)
			Expect(err).NotTo(HaveOccurred())

			tmpDir := os.TempDir()
			childItems, err := ioutil.ReadDir(tmpDir)
			Expect(err).NotTo(HaveOccurred())

			for _, item := range childItems {
				if item.IsDir() && strings.HasPrefix(filepath.Base(item.Name()), "pool-resource") {
					fmt.Printf("Cleaning up temporary pool resource %s\n", item.Name())
					_ = os.RemoveAll(item.Name())
				}
			}
		}
	}
}, func() {
	Expect(os.RemoveAll(stembuildExecutable)).To(Succeed())
})
