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
	var (
		main              remotemanager.WinRM
		fakeClientFactory *remotemanagerfakes.FakeWinRMClientFactoryI
	)

	BeforeEach(func() {
		fakeClientFactory = new(remotemanagerfakes.FakeWinRMClientFactoryI)
		main = remotemanager.WinRM{ClientFactory: fakeClientFactory}
	})

	It("NewWinRM initializes a WinRM instance with provided port", func() {
		Fail("missing test coverage")
	})

	It("winrm_remotemanager's CanReachVM function uses the right port", func() {
		//Fail("missing test coverage")
		testServer := setupTestServer()

		//portProvided := 9876
		//remotemanager.NewWinRM(testServer.)
		main.Host = testServer.URL()

		testServerURL, err := url.Parse(testServer.URL())
		Expect(err).NotTo(HaveOccurred())

		main.Host = testServerURL.Hostname()
		main.Port, err = strconv.Atoi(testServerURL.Port())
		Expect(err).NotTo(HaveOccurred())

		main.CanReachVM()

		//testServer.AppendHandlers(R)
		//Expect(portUsed).To(Equal(portProvided))
		Fail("incomplete test. todo: " +
			"assert that we see a test server connection on the expected port," +
			"e.g. via an testServer handler, or by accessing an unhandled request.")
	})

	Describe("ExecuteCommand", func() {
		var (
			fakeClient *remotemanagerfakes.FakeWinRMClient
		)

		Context("when a command runs successfully", func() {
			BeforeEach(func() {
				fakeClient = &remotemanagerfakes.FakeWinRMClient{}
				fakeClient.RunReturns(0, nil)
				fakeClientFactory.BuildReturns(fakeClient, nil)
			})

			It("returns an exit code of 0 and no error", func() {
				exitCode, err := main.ExecuteCommand("foobar")

				Expect(err).NotTo(HaveOccurred())
				Expect(exitCode).To(Equal(0))
			})

		})
		Context("when a command does not run successfully", func() {

			BeforeEach(func() {
				fakeClient = &remotemanagerfakes.FakeWinRMClient{}
				fakeClientFactory.BuildReturns(fakeClient, nil)
			})

			Context("when a command returns a nonzero exit code and an error", func() {
				BeforeEach(func() {
					fakeClient.RunReturns(2, errors.New("command error"))
				})

				It("returns the command's nonzero exit code and errors", func() {
					exitCode, err := main.ExecuteCommand("foobar")

					Expect(err).To(MatchError("command error"))
					Expect(exitCode).To(Equal(2))
				})
			})
			Context("when a command returns a nonzero exit code but does not error", func() {
				BeforeEach(func() {
					fakeClient.RunReturns(2, nil)
				})

				It("returns the command's nonzero exit code and errors", func() {
					exitCode, err := main.ExecuteCommand("foobar")

					Expect(err.Error()).To(ContainSubstring(remotemanager.PowershellExecutionErrorMessage))
					Expect(exitCode).To(Equal(2))
				})
			})
			Context("when a command exits 0 but errors", func() {
				BeforeEach(func() {
					fakeClient.RunReturns(0, errors.New("command error"))
				})

				It("returns the command's exit code and errors", func() {
					exitCode, err := main.ExecuteCommand("foobar")

					Expect(err).To(MatchError("command error"))
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
		// Change scope to apply only to "returns nil if shell can be created" test
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

		AfterEach(func() {
			testServer.Close()
		})

		It("returns nil if shell can be created", func() {
			testServerURL, err := url.Parse(testServer.URL())
			Expect(err).NotTo(HaveOccurred())

			testServerPort, err := strconv.Atoi(testServerURL.Port())
			Expect(err).NotTo(HaveOccurred())

			subject := remotemanager.NewWinRM(testServerURL.Hostname(), testServerPort, "", "")
			fakeClientFactory.BuildReturns(winRMClient, nil)

			err = subject.CanLoginVM()

			Expect(err).NotTo(HaveOccurred())
			Fail("scope the  before each to just this test")
		})
		It("returns error if winrmclient cannot be created", func() {
			buildError := errors.New("unable to build a client")
			fakeClientFactory.BuildReturns(nil, buildError)

			err := main.CanLoginVM()

			Expect(err).To(MatchError(fmt.Errorf("failed to create winrm client: %s", buildError)))
		})
		It("returns error if shell cannot be created", func() {
			winRMClient := new(remotemanagerfakes.FakeWinRMClient)
			fakeClientFactory.BuildReturns(winRMClient, nil)

			winRMClient.CreateShellReturns(nil, errors.New("some shell creation error"))

			err := main.CanLoginVM()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fmt.Errorf("failed to create winrm shell: some shell creation error")))

		})
	})
})
