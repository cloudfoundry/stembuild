package vcenter_client_factory_test

import (
	"context"
	"errors"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25"

	vcenterclientfactory "github.com/cloudfoundry/stembuild/iaas_cli/iaas_clients/factory"
	"github.com/cloudfoundry/stembuild/iaas_cli/iaas_clients/factory/factoryfakes"
	"github.com/cloudfoundry/stembuild/iaas_cli/iaas_clients/vcenter_manager"
)

var _ = Describe("VcenterManagerFactory", func() {

	var (
		managerFactory *vcenterclientfactory.ManagerFactory
	)

	BeforeEach(func() {
		managerFactory = &vcenterclientfactory.ManagerFactory{}
	})

	Context("VCenterManager", func() {
		It("returns a vcenter manager", func() {

			fakeVimClient := &vim25.Client{}
			fakeClientCreator := &factoryfakes.FakeVim25ClientCreator{}

			fakeClientCreator.NewClientReturns(fakeVimClient, nil)

			fakeFinder := &find.Finder{}
			fakeFinderCreator := &factoryfakes.FakeFinderCreator{}
			fakeFinderCreator.NewFinderReturns(fakeFinder)

			managerFactory.SetConfig(vcenterclientfactory.FactoryConfig{
				VCenterServer:  "example.com",
				Username:       "user",
				Password:       "pass",
				ClientCreator:  fakeClientCreator,
				FinderCreator:  fakeFinderCreator,
				RootCACertPath: "",
			})

			manager, err := managerFactory.VCenterManager(context.TODO())
			Expect(err).NotTo(HaveOccurred())

			Expect(manager).To(BeAssignableToTypeOf(&vcenter_manager.VCenterManager{}))

		})

		It("returns an error if the vcenter server cannot be parsed", func() {
			fakeClientCreator := &factoryfakes.FakeVim25ClientCreator{}

			managerFactory.SetConfig(vcenterclientfactory.FactoryConfig{
				VCenterServer: " :", // make soap.ParseURL fail with
				Username:      "user",
				Password:      "pass",
				ClientCreator: fakeClientCreator,
			})

			_, err := managerFactory.VCenterManager(context.TODO())
			var parseErr *url.Error
			ok := errors.As(err, &parseErr)
			Expect(ok).To(BeTrue())
			Expect(parseErr.Op).To(Equal("parse"))
		})

		It("returns an error if a vim25 client cannot be created", func() {

			clientErr := errors.New("can't make a client")
			fakeClientCreator := &factoryfakes.FakeVim25ClientCreator{}
			fakeClientCreator.NewClientReturns(nil, clientErr)

			managerFactory.SetConfig(vcenterclientfactory.FactoryConfig{
				VCenterServer: "example.com",
				Username:      "user",
				Password:      "pass",
				ClientCreator: fakeClientCreator,
			})

			_, err := managerFactory.VCenterManager(context.TODO())
			Expect(err).To(MatchError(clientErr))
		})
	})
})
