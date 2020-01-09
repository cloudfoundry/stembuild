package vmconstruct_factory

import (
	"github.com/cloudfoundry-incubator/stembuild/commandparser/commandparserfakes"
	"github.com/cloudfoundry-incubator/stembuild/construct"
	"github.com/cloudfoundry-incubator/stembuild/construct/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("Factory", func() {
	Describe("GetVMPreparer", func() {
		var (
			factory *VMConstructFactory
		)

		BeforeEach(func() {
			factory = &VMConstructFactory{}
		})

		It("should return a VMPreparer", func() {
			fakeVCenterManager := &commandparserfakes.FakeVCenterManager{}

			sourceConfig := config.SourceConfig{
				GuestVmIp:       "vmIP",
				GuestVMUsername: "vmUser",
				GuestVMPassword: "vmPwd",
				VCenterUrl:      "vCenterUrl",
				VCenterUsername: "vCenterUser",
				VCenterPassword: "vCenterPwd",
				VmInventoryPath: "some-vm-inventory-path",
			}

			vmPreparer, err := factory.VMPreparer(sourceConfig, fakeVCenterManager)
			Expect(err).ToNot(HaveOccurred())
			Expect(vmPreparer).To(BeAssignableToTypeOf(&construct.VMConstruct{}))
		})

		It("should return a login error when login incorrect to VCenter", func() {
			// setup
			fakeVCenterManager := &commandparserfakes.FakeVCenterManager{}
			loginFailure := errors.New("could not log in")
			fakeVCenterManager.LoginReturns(loginFailure)
			sourceConfig := config.SourceConfig{}

			vmPreparer, err := factory.VMPreparer(sourceConfig, fakeVCenterManager)

			Expect(vmPreparer).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Cannot complete login due to an incorrect vCenter user name or password"))
			Expect(err.Error()).To(ContainSubstring(loginFailure.Error()))
		})
	})
})
