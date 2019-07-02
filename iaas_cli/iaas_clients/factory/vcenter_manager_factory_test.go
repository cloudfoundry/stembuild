package vcenter_client_factory_test

import (
	"context"
	"errors"

	"github.com/vmware/govmomi/find"

	"github.com/cloudfoundry-incubator/stembuild/iaas_cli/iaas_clients/vcenter_manager"

	"github.com/vmware/govmomi/vim25"

	"github.com/cloudfoundry-incubator/stembuild/iaas_cli/iaas_clients/factory"
	"github.com/cloudfoundry-incubator/stembuild/iaas_cli/iaas_clients/factory/factoryfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VcenterManagerFactory", func() {

	Context("VCenterManager", func() {
		It("returns a vcenter manager", func() {

			fakeVimClient := &vim25.Client{}
			fakeClientCreator := &factoryfakes.FakeVim25ClientCreator{}

			fakeClientCreator.NewClientReturns(fakeVimClient, nil)

			fakeFinder := &find.Finder{}
			fakeFinderCreator := &factoryfakes.FakeFinderCreator{}
			fakeFinderCreator.NewFinderReturns(fakeFinder)

			f := vcenter_client_factory.ManagerFactory{
				VCenterServer:  "example.com",
				Username:       "user",
				Password:       "pass",
				ClientCreator:  fakeClientCreator,
				FinderCreator:  fakeFinderCreator,
				RootCACertPath: "",
			}

			manager, err := f.VCenterManager(context.TODO())
			Expect(err).NotTo(HaveOccurred())

			Expect(manager).To(BeAssignableToTypeOf(&vcenter_manager.VCenterManager{}))

		})

		It("returns an error if the vcenter server cannot be parsed", func() {

			parseErr := errors.New("net/url: invalid control character in URL")
			fakeClientCreator := &factoryfakes.FakeVim25ClientCreator{}

			f := vcenter_client_factory.ManagerFactory{
				VCenterServer: string(127), // make soap.ParseURL fail with
				Username:      "user",
				Password:      "pass",
				ClientCreator: fakeClientCreator,
			}

			_, err := f.VCenterManager(context.TODO())
			Expect(err.Error()).To(ContainSubstring(parseErr.Error()))

		})

		It("returns an error if a vim25 client cannot be created", func() {

			clientErr := errors.New("can't make a client")
			fakeClientCreator := &factoryfakes.FakeVim25ClientCreator{}
			fakeClientCreator.NewClientReturns(nil, clientErr)

			f := vcenter_client_factory.ManagerFactory{
				VCenterServer: "example.com",
				Username:      "user",
				Password:      "pass",
				ClientCreator: fakeClientCreator,
			}

			_, err := f.VCenterManager(context.TODO())
			Expect(err).To(MatchError(clientErr))
		})
	})
})
