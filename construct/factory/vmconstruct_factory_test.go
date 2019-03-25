package vmconstruct_factory

import (
	"github.com/cloudfoundry-incubator/stembuild/construct"
	"github.com/cloudfoundry-incubator/stembuild/construct/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			sourceConfig := config.SourceConfig{
				GuestVmIp:       "vmIP",
				GuestVMUsername: "vmUser",
				GuestVMPassword: "vmPwd",
				VCenterUrl:      "vCenterUrl",
				VCenterUsername: "vCenterUser",
				VCenterPassword: "vCenterPwd",
				VmInventoryPath: "some-vm-inventory-path",
			}

			vmPreparer := factory.VMPreparer(sourceConfig)
			Expect(vmPreparer).To(BeAssignableToTypeOf(&construct.VMConstruct{}))
		})
	})
})
