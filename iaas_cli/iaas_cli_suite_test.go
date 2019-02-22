package iaas_cli_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIaasCli(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IaasCli Suite")
}

var targetVMPath string

var _ = BeforeSuite(func() {
	vmFolder := os.Getenv("VM_FOLDER")
	Expect(vmFolder).NotTo(Equal(""), "VM_FOLDER is required")
	vmName := os.Getenv("PACKAGE_TEST_BASE_VM_NAME")
	Expect(vmName).NotTo(Equal(""), "PACKAGE_TEST_BASE_VM_NAME is required")

	targetVMPath = fmt.Sprintf("%s/%s", vmFolder, vmName)
})
