package iaas_clients

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/stembuild/iaas_cli/iaas_clients/factory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	VcenterUrl      = "VCENTER_URL"
	VcenterUsername = "VCENTER_USERNAME"
	VcenterPassword = "VCENTER_PASSWORD"
	TestVmPath      = "CONTRACT_TEST_VM_PATH"
	TestVmPassword  = "CONTRACT_TEST_VM_PASSWORD"
	TestVmUsername  = "CONTRACT_TEST_VM_USERNAME"
)

func envMustExist(variableName string) string {
	result := os.Getenv(variableName)
	if result == "" {
		Fail(fmt.Sprintf("%s must be set", variableName))
	}

	return result
}

//TODO: test for certs
//TODO: create test vm dynamically from a base
var _ = Describe("VcenterClient", func() {
	Describe("StartProgram", func() {

		It("Starts a program and returns its exit code", func() {

			managerFactory := vcenter_client_factory.ManagerFactory{
				VCenterServer:      envMustExist(VcenterUrl),
				Username:           envMustExist(VcenterUsername),
				Password:           envMustExist(VcenterPassword),
				InsecureConnection: true,
				ClientCreator:      &vcenter_client_factory.ClientCreator{},
				FinderCreator:      &vcenter_client_factory.GovmomiFinderCreator{},
			}

			ctx := context.TODO()
			vCenterManager, err := managerFactory.VCenterManager(ctx)
			Expect(err).ToNot(HaveOccurred())

			vCenterManager.Login(ctx)

			vm, err := vCenterManager.FindVM(ctx, envMustExist(TestVmPath))
			Expect(err).ToNot(HaveOccurred())

			opsManager := vCenterManager.OperationsManager(ctx, vm)

			guestManager, err := vCenterManager.GuestManager(ctx, opsManager, envMustExist(TestVmUsername), envMustExist(TestVmPassword))
			Expect(err).ToNot(HaveOccurred())

			powershell := "C:\\Windows\\System32\\WindowsPowerShell\\V1.0\\powershell.exe"
			pid, err := guestManager.StartProgramInGuest(ctx, powershell, "Exit 59")
			Expect(err).ToNot(HaveOccurred())

			exitCode, err := guestManager.ExitCodeForProgramInGuest(ctx, pid)
			Expect(err).ToNot(HaveOccurred())
			Expect(exitCode).To(Equal(int32(59)))
		})
	})
})
