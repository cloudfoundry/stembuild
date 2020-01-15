package construct_test

import (
	"errors"
	//"github.com/cloudfoundry-incubator/stembuild/construct"
	"github.com/cloudfoundry-incubator/stembuild/construct"
	"github.com/cloudfoundry-incubator/stembuild/remotemanager/remotemanagerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//. "github.com/onsi/gomega"
)

var _ = Describe("VMConnectionValidator", func() {
	var (
		validator         *construct.WinRMConnectionValidator
		fakeRemoteManager *remotemanagerfakes.FakeRemoteManager
	)

	BeforeEach(func() {
		fakeRemoteManager = &remotemanagerfakes.FakeRemoteManager{}

		validator = &construct.WinRMConnectionValidator{
			RemoteManager: fakeRemoteManager,
		}
	})

	Describe("Validate connection to the VM", func() {
		It("can reach VM and can login to VM", func() {
			err := validator.Validate()

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeRemoteManager.CanReachVMCallCount()).To(Equal(1))
			Expect(fakeRemoteManager.CanLoginVMCallCount()).To(Equal(1))
		})

		It("return an error when it cannot reach the VM", func() {
			fakeRemoteManager.CanReachVMReturns(errors.New("could not reach vm"))
			err := validator.Validate()

			Expect(err).To(HaveOccurred())
			Expect(fakeRemoteManager.CanReachVMCallCount()).To(Equal(1))
			Expect(fakeRemoteManager.CanLoginVMCallCount()).To(Equal(0))
		})

		It("return an error when it cannot log into the VM", func() {
			invalidLoginError := errors.New("Cannot complete login due to an incorrect VM user name or password")
			fakeRemoteManager.CanLoginVMReturns(errors.New("login error"))

			err := validator.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(errors.New("login error").Error()))
			Expect(err.Error()).To(ContainSubstring(invalidLoginError.Error()))

			Expect(fakeRemoteManager.CanReachVMCallCount()).To(Equal(1))
			Expect(fakeRemoteManager.CanLoginVMCallCount()).To(Equal(1))
		})
	})

})
