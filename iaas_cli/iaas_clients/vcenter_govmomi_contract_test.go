package iaas_clients

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/stembuild/iaas_cli/iaas_clients/factory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	VcenterUrl      = "VCENTER_URL"
	VcenterUsername = "VCENTER_USERNAME"
	VcenterPassword = "VCENTER_PASSWORD"
	VcenterCACert   = "VCENTER_CA_CERT"
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

//TODO: create test vm dynamically from a base
var _ = Describe("VcenterClient", func() {
	FDescribe("StartProgram", func() {

		var managerFactory *vcenter_client_factory.ManagerFactory

		ExpectProgramToStartAndExitSuccessfully := func() {

			ctx := context.TODO()

			vCenterManager, err := managerFactory.VCenterManager(ctx)
			ExpectWithOffset(1, err).ToNot(HaveOccurred())

			err = vCenterManager.Login(ctx)
			ExpectWithOffset(1, err).ToNot(HaveOccurred())

			vm, err := vCenterManager.FindVM(ctx, envMustExist(TestVmPath))
			ExpectWithOffset(1, err).ToNot(HaveOccurred())

			opsManager := vCenterManager.OperationsManager(ctx, vm)
			guestManager, err := vCenterManager.GuestManager(ctx, opsManager, envMustExist(TestVmUsername), envMustExist(TestVmPassword))
			ExpectWithOffset(1, err).ToNot(HaveOccurred())

			powershell := "C:\\Windows\\System32\\WindowsPowerShell\\V1.0\\powershell.exe"
			pid, err := guestManager.StartProgramInGuest(ctx, powershell, "Exit 59")
			ExpectWithOffset(1, err).ToNot(HaveOccurred())

			exitCode, err := guestManager.ExitCodeForProgramInGuest(ctx, pid)
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, exitCode).To(Equal(int32(59)))
		}

		BeforeEach(func() {
			managerFactory = &vcenter_client_factory.ManagerFactory{
				VCenterServer: envMustExist(VcenterUrl),
				Username:      envMustExist(VcenterUsername),
				Password:      envMustExist(VcenterPassword),
				ClientCreator: &vcenter_client_factory.ClientCreator{},
				FinderCreator: &vcenter_client_factory.GovmomiFinderCreator{},
			}
		})

		AfterEach(func() {
			managerFactory = nil
		})

		Context("Use root cert implicitly", func() {
			It("Starts a program and returns its exit code", func() {

				managerFactory.RootCACertPath = ""
				ExpectProgramToStartAndExitSuccessfully()
			})
		})

		Context("A factory is given a proper CA cert", func() {

			It("Starts a program and returns its exit code", func() {

				cert := os.Getenv(VcenterCACert)
				if cert == "" {
					Skip(fmt.Sprintf("export VCENTER_CA_CERT=<a valid ca cert> to run this test"))
				}

				tmpDir, err := ioutil.TempDir("", "vcenter-client-contract-tests")
				defer os.RemoveAll(tmpDir)
				Expect(err).ToNot(HaveOccurred())
				f, err := ioutil.TempFile(tmpDir, "valid-cert")
				Expect(err).ToNot(HaveOccurred())

				_, err = f.WriteString(cert)
				Expect(err).ToNot(HaveOccurred())

				err = f.Close()
				Expect(err).ToNot(HaveOccurred())

				managerFactory.RootCACertPath = f.Name()

				ExpectProgramToStartAndExitSuccessfully()
			})

		})

		Context("A factory is given an improper CA cert", func() {

			It("fails to create a vcenter manager", func() {

				workingDir, err := os.Getwd()
				Expect(err).NotTo(HaveOccurred())
				fakeCertPath := filepath.Join(workingDir, "fixtures", "fakecert")

				managerFactory.RootCACertPath = fakeCertPath

				ctx := context.TODO()
				_, err = managerFactory.VCenterManager(ctx)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be used as a trusted CA certificate"))

			})
		})

	})
})
