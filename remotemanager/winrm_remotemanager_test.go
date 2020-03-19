package remotemanager_test

import (
	"errors"
	"github.com/cloudfoundry-incubator/stembuild/remotemanager"
	"github.com/cloudfoundry-incubator/stembuild/remotemanager/remotemanagerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WinRM RemoteManager", func() {
	Describe("ExecuteCommand", func() {
		var (
			fakeClientFactory *remotemanagerfakes.FakeWinRMClientFactoryI
			fakeClient        *remotemanagerfakes.FakeWinRMClient
		)

		Context("when a command runs successfully", func() {
			BeforeEach(func() {
				fakeClient = &remotemanagerfakes.FakeWinRMClient{}
				fakeClient.RunReturns(0, nil)
				fakeClientFactory = &remotemanagerfakes.FakeWinRMClientFactoryI{}
				fakeClientFactory.BuildReturns(fakeClient, nil)
			})

			It("returns an exit code of 0 and no error", func() {
				remoteManager := remotemanager.NewWinRM("foo", "bar", "baz", fakeClientFactory)
				exitCode, err := remoteManager.ExecuteCommand("foobar")

				Expect(err).NotTo(HaveOccurred())
				Expect(exitCode).To(Equal(0))
			})

		})
		Context("when a command does not run successfully", func() {

			BeforeEach(func() {
				fakeClient = &remotemanagerfakes.FakeWinRMClient{}
				fakeClientFactory = &remotemanagerfakes.FakeWinRMClientFactoryI{}
				fakeClientFactory.BuildReturns(fakeClient, nil)
			})

			Context("when a command returns a nonzero exit code and an error", func() {
				BeforeEach(func() {
					fakeClient.RunReturns(2, errors.New("command error"))
				})

				It("returns the command's nonzero exit code and errors", func() {
					remoteManager := remotemanager.NewWinRM("foo", "bar", "baz", fakeClientFactory)
					exitCode, err := remoteManager.ExecuteCommand("foobar")

					Expect(err).To(HaveOccurred())
					Expect(exitCode).To(Equal(2))
				})
			})
			Context("when a command returns a nonzero exit code but does not error", func() {
				BeforeEach(func() {
					fakeClient.RunReturns(2, nil)
				})

				It("returns the command's nonzero exit code and errors", func() {
					remoteManager := remotemanager.NewWinRM("foo", "bar", "baz", fakeClientFactory)
					exitCode, err := remoteManager.ExecuteCommand("foobar")

					Expect(err).To(HaveOccurred())
					Expect(exitCode).To(Equal(2))
				})
			})
			Context("when a command exits 0 but errors", func() {
				BeforeEach(func() {
					fakeClient.RunReturns(0, errors.New("command error"))
				})

				It("returns the command's exit code and errors", func() {
					remoteManager := remotemanager.NewWinRM("foo", "bar", "baz", fakeClientFactory)
					exitCode, err := remoteManager.ExecuteCommand("foobar")

					Expect(err).To(HaveOccurred())
					Expect(exitCode).To(Equal(0))
				})
			})
		})

	})
})
