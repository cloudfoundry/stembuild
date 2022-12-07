package vcenter_manager_test

import (
	"context"
	"errors"
	"runtime"

	"github.com/cloudfoundry/stembuild/iaas_cli/iaas_clients/factory"

	"github.com/cloudfoundry/stembuild/iaas_cli/iaas_clients/guest_manager"
	"github.com/vmware/govmomi/guest"
	"github.com/vmware/govmomi/object"

	"github.com/cloudfoundry/stembuild/iaas_cli/iaas_clients/vcenter_manager"
	"github.com/cloudfoundry/stembuild/iaas_cli/iaas_clients/vcenter_manager/vcenter_managerfakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vmware/govmomi/vim25"
)

var _ = Describe("VcenterManager", func() {

	var (
		fakeGovmomiClient vcenter_managerfakes.FakeGovmomiClient
		fakeVimClient     vim25.Client
		fakeFinder        vcenter_managerfakes.FakeFinder
	)

	BeforeEach(func() {
		fakeGovmomiClient = vcenter_managerfakes.FakeGovmomiClient{}
		fakeVimClient = vim25.Client{}
	})

	Context("Login", func() {

		It("logs the user into vcenter", func() {

			fakeGovmomiClient.LoginReturns(nil)

			vcManager, err := vcenter_manager.NewVCenterManager(&fakeGovmomiClient, &fakeVimClient, &fakeFinder, "user", "pass")
			Expect(err).ToNot(HaveOccurred())

			err = vcManager.Login(context.TODO())

			_, u := fakeGovmomiClient.LoginArgsForCall(0)
			pass, _ := u.Password()
			Expect(u.Username()).To(Equal("user"))
			Expect(pass).To(Equal("pass"))

			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error if the client encounters one", func() {

			loginErr := errors.New("bummer dude")
			fakeGovmomiClient.LoginReturns(loginErr)

			vcManager, err := vcenter_manager.NewVCenterManager(&fakeGovmomiClient, &fakeVimClient, &fakeFinder, "user", "pass")
			Expect(err).ToNot(HaveOccurred())

			err = vcManager.Login(context.TODO())
			Expect(err).To(MatchError(loginErr))
		})
	})

	Context("FindVM", func() {

		It("searches for the specified vm", func() {

			fakeVM := &object.VirtualMachine{}
			fakeFinder.VirtualMachineReturns(fakeVM, nil)

			vcManager, err := vcenter_manager.NewVCenterManager(&fakeGovmomiClient, &fakeVimClient, &fakeFinder, "user", "pass")
			Expect(err).ToNot(HaveOccurred())

			vm, err := vcManager.FindVM(context.TODO(), "/path/to/some/vm")
			Expect(err).ToNot(HaveOccurred())
			_, path := fakeFinder.VirtualMachineArgsForCall(0)
			Expect(path).To(Equal("/path/to/some/vm"))

			Expect(vm).To(BeEquivalentTo(fakeVM))
		})

		It("returns an error if the finder does", func() {

			findErr := errors.New("can't find it, friend.")
			fakeFinder.VirtualMachineReturns(nil, findErr)

			vcManager, err := vcenter_manager.NewVCenterManager(&fakeGovmomiClient, &fakeVimClient, &fakeFinder, "user", "pass")
			Expect(err).ToNot(HaveOccurred())

			_, err = vcManager.FindVM(context.TODO(), "/path/to/some/vm")
			Expect(err).To(MatchError(findErr))

		})
	})

	Context("GuestManager", func() {

		It("searches for the specified vm", func() {

			fakeProcManager := &guest.ProcessManager{}
			fakeOpsManager := &vcenter_managerfakes.FakeOpsManager{}
			fakeOpsManager.ProcessManagerReturns(fakeProcManager, nil)

			vcManager, err := vcenter_manager.NewVCenterManager(&fakeGovmomiClient, &fakeVimClient, &fakeFinder, "user", "pass")
			Expect(err).ToNot(HaveOccurred())

			gm, err := vcManager.GuestManager(context.TODO(), fakeOpsManager, "guestUser", "guestPass")
			Expect(err).ToNot(HaveOccurred())

			Expect(gm).To(BeAssignableToTypeOf(&guest_manager.GuestManager{}))
		})

		It("returns an error if the finder does", func() {

			guestErr := errors.New("not today, junior")
			fakeOpsManager := &vcenter_managerfakes.FakeOpsManager{}
			fakeOpsManager.ProcessManagerReturns(nil, guestErr)

			vcManager, err := vcenter_manager.NewVCenterManager(&fakeGovmomiClient, &fakeVimClient, &fakeFinder, "user", "pass")
			Expect(err).ToNot(HaveOccurred())

			_, err = vcManager.GuestManager(context.TODO(), fakeOpsManager, "guestUser", "guestPass")
			Expect(err).To(MatchError(guestErr))

		})
	})

	Context("running against vcsim server", func() {
		Context("CloneVM", func() {
			It("clones a vm", func() {
				if runtime.GOOS == "windows" {
					Skip("windows cannot run a vcsim server")
				}

				inventoryPath := "/DC0/vm/DC0_H0_VM0"
				clonePath := "/DC0/vm/DC0_H0_VM0_NewClone"

				factoryConfig := &vcenter_client_factory.FactoryConfig{
					VCenterServer:  "https://user:pass@127.0.0.1:8989/sdk",
					Username:       "user",
					Password:       "pass",
					ClientCreator:  &vcenter_client_factory.ClientCreator{},
					FinderCreator:  &vcenter_client_factory.GovmomiFinderCreator{},
					RootCACertPath: CertPath,
				}

				managerFactory := &vcenter_client_factory.ManagerFactory{
					*factoryConfig,
				}

				ctx := context.TODO()

				vCenterManager, err := managerFactory.VCenterManager(ctx)
				Expect(err).ToNot(HaveOccurred())

				err = vCenterManager.Login(ctx)
				Expect(err).ToNot(HaveOccurred())

				vmToClone, err := vCenterManager.FindVM(ctx, inventoryPath)
				Expect(err).ToNot(HaveOccurred())

				err = vCenterManager.CloneVM(ctx, vmToClone, clonePath)
				Expect(err).ToNot(HaveOccurred())

				_, err = vCenterManager.FindVM(ctx, clonePath)
				Expect(err).ToNot(HaveOccurred())

			})
		})
	})
})
