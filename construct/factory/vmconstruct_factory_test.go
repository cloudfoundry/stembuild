package vmconstruct_factory

import (
	"github.com/cloudfoundry-incubator/stembuild/construct"
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
			vmPreparer := factory.VMPreparer("0.0.0.0", "pivotal", "password")
			Expect(vmPreparer).To(BeAssignableToTypeOf(&construct.VMConstruct{}))
		})
	})
})
