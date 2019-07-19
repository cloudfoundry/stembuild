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

func envMustExist(variableName string) string {
	result := os.Getenv(variableName)
	if result == "" {
		Fail(fmt.Sprintf("%s must be set", variableName))
	}

	return result
}

var _ = Describe("VcenterClient", func() {
	Describe("StartProgram", func() {

		var (
			managerFactory *vcenter_client_factory.ManagerFactory
			factoryConfig  *vcenter_client_factory.FactoryConfig
		)
		ExpectProgramToStartAndExitSuccessfully := func() {

			ctx := context.TODO()

			vCenterManager, err := managerFactory.VCenterManager(ctx)
			Expect(err).ToNot(HaveOccurred())

			err = vCenterManager.Login(ctx)
			Expect(err).ToNot(HaveOccurred())

			vm, err := vCenterManager.FindVM(ctx, TestVmPath)
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
		}

		BeforeEach(func() {

			factoryConfig = &vcenter_client_factory.FactoryConfig{
				VCenterServer: envMustExist(VcenterUrl),
				Username:      envMustExist(VcenterUsername),
				Password:      envMustExist(VcenterPassword),
				ClientCreator: &vcenter_client_factory.ClientCreator{},
				FinderCreator: &vcenter_client_factory.GovmomiFinderCreator{},
			}

			managerFactory = &vcenter_client_factory.ManagerFactory{}
		})

		AfterEach(func() {
			managerFactory = nil
		})

		Context("Use root cert implicitly", func() {
			It("Starts a program and returns its exit code", func() {

				factoryConfig.RootCACertPath = ""
				managerFactory.SetConfig(*factoryConfig)
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

				factoryConfig.RootCACertPath = f.Name()
				managerFactory.SetConfig(*factoryConfig)

				ExpectProgramToStartAndExitSuccessfully()
			})

		})

		Context("A factory is given an improper CA cert", func() {

			It("fails to create a vcenter manager", func() {

				workingDir, err := os.Getwd()
				Expect(err).NotTo(HaveOccurred())
				fakeCertPath := filepath.Join(workingDir, "fixtures", "fakecert")

				factoryConfig.RootCACertPath = fakeCertPath
				managerFactory.SetConfig(*factoryConfig)

				ctx := context.TODO()
				_, err = managerFactory.VCenterManager(ctx)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be used as a trusted CA certificate"))

			})
		})

	})
})
