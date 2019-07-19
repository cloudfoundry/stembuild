package iaas_clients

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/vmware/govmomi/object"

	"github.com/cloudfoundry-incubator/stembuild/iaas_cli/iaas_clients/factory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIaasClients(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IaasClients Suite")
}

const (
	VcenterUrl      = "VCENTER_BASE_URL"
	VcenterUsername = "VCENTER_USERNAME"
	VcenterPassword = "VCENTER_PASSWORD"
	VcenterCACert   = "VCENTER_CA_CERT"
	BaseVmPath      = "CONTRACT_TEST_VM_PATH"
	TestVmPassword  = "CONTRACT_TEST_VM_PASSWORD"
	TestVmUsername  = "CONTRACT_TEST_VM_USERNAME"
)

var TestVmPath string
var VM *object.VirtualMachine
var CTX context.Context
var _ = BeforeSuite(func() {

	managerFactory := &vcenter_client_factory.ManagerFactory{Config: vcenter_client_factory.FactoryConfig{
		VCenterServer: envMustExist(VcenterUrl),
		Username:      envMustExist(VcenterUsername),
		Password:      envMustExist(VcenterPassword),
		ClientCreator: &vcenter_client_factory.ClientCreator{},
		FinderCreator: &vcenter_client_factory.GovmomiFinderCreator{},
	},
	}

	CTX = context.TODO()

	vCenterManager, err := managerFactory.VCenterManager(CTX)
	Expect(err).ToNot(HaveOccurred())

	err = vCenterManager.Login(CTX)
	Expect(err).ToNot(HaveOccurred())

	vmToClone, err := vCenterManager.FindVM(CTX, envMustExist(BaseVmPath))
	Expect(err).ToNot(HaveOccurred())

	TestVmPath = envMustExist(BaseVmPath) + fmt.Sprintf("%s", uuid.New())[0:8]

	err = vCenterManager.CloneVM(CTX, vmToClone, TestVmPath)
	Expect(err).ToNot(HaveOccurred())

	time.Sleep(30 * time.Second)

	VM, err = vCenterManager.FindVM(CTX, TestVmPath)
	Expect(err).ToNot(HaveOccurred())

})

var _ = AfterSuite(func() {

	if VM != nil {
		task, err := VM.PowerOff(CTX)
		Expect(err).ToNot(HaveOccurred())
		err = task.Wait(CTX)
		Expect(err).ToNot(HaveOccurred())

		task, err = VM.Destroy(CTX)
		Expect(err).ToNot(HaveOccurred())
		err = task.Wait(CTX)
		Expect(err).ToNot(HaveOccurred())
	}
})
