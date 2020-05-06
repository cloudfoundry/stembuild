package remotemanager_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cloudfoundry-incubator/stembuild/remotemanager"
	"github.com/cloudfoundry-incubator/stembuild/remotemanager/remotemanagerfakes"
	"github.com/masterzen/winrm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

////go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . FakeShell
//type FakeShell interface {
//	 Close() error
//}

func setupTestServer() *Server {
	server := NewServer()

	// Suppresses ginkgo server logs
	server.HTTPTestServer.Config.ErrorLog = log.New(&bytes.Buffer{}, "", 0)

	response := `
<s:Envelope xmlns:s="https://www.w3.org/2003/05/soap-envelope" 
            xmlns:a="https://schemas.xmlsoap.org/ws/2004/08/addressing" 
            xmlns:w="https://schemas.dmtf.org/wbem/wsman/1/wsman.xsd">
  <s:Header>
    <w:ShellId>153600</w:ShellId>
  </s:Header>
  <s:Body/> 
</s:Envelope>
` // Looks for: //w:Selector[@Name='ShellId']
	server.AppendHandlers(
		RespondWith(http.StatusOK, response, http.Header{
			"Content-Type": {"application/soap+xml;charset=UTF-8"},
		}),
	)
	server.AllowUnhandledRequests = true
	return server
}

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
	Describe("CanLoginVM", func() {
		var (
			testServer  *Server
			winRMClient *winrm.Client
		)

		BeforeEach(func() {
			testServer = setupTestServer()

			testServerURL, err := url.Parse(testServer.URL())
			Expect(err).NotTo(HaveOccurred())
			port, err := strconv.Atoi(testServerURL.Port())
			Expect(err).NotTo(HaveOccurred())

			endpoint := &winrm.Endpoint{
				Host: testServerURL.Hostname(),
				Port: port,
			}
			winRMClient, err = winrm.NewClient(endpoint, "", "")
			Expect(err).NotTo(HaveOccurred())
		})

		var _ = AfterEach(func() {
			testServer.Close()
		})

		It("returns nil if shell can be created", func() {
			winRMClientFactory := new(remotemanagerfakes.FakeWinRMClientFactoryI)
			winRMClientFactory.BuildReturns(winRMClient, nil)

			remotemanager := remotemanager.NewWinRM("some-host", "some-user", "some-pass", winRMClientFactory)

			err := remotemanager.CanLoginVM()

			Expect(err).NotTo(HaveOccurred())
		})
		It("returns error if winrmclient cannot be created", func() {
			winRMClientFactory := new(remotemanagerfakes.FakeWinRMClientFactoryI)
			buildError := errors.New("unable to build a client")
			winRMClientFactory.BuildReturns(nil, buildError)

			remotemanager := remotemanager.NewWinRM("some-host", "some-user", "some-pass", winRMClientFactory)
			err := remotemanager.CanLoginVM()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fmt.Errorf("failed to create winrm client: %s", buildError)))
		})
		It("returns error if shell cannot be created", func() {
			winRMClientFactory := new(remotemanagerfakes.FakeWinRMClientFactoryI)
			winRMClient := new(remotemanagerfakes.FakeWinRMClient)
			winRMClientFactory.BuildReturns(winRMClient, nil)

			winRMClient.CreateShellReturns(nil, errors.New("some shell creation error"))
			remotemanager := remotemanager.NewWinRM("some-host", "some-user", "some-pass", winRMClientFactory)

			err := remotemanager.CanLoginVM()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fmt.Errorf("failed to create winrm shell: some shell creation error")))

		})
	})
})
