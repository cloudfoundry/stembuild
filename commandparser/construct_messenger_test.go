package commandparser_test

import (
	"errors"
	"github.com/cloudfoundry-incubator/stembuild/commandparser"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("ConstructMessenger", func() {
	var (
		cm commandparser.ConstructCmdMessenger
		g  *Buffer
	)

	BeforeEach(func() {
		g = NewBuffer()
		cm = commandparser.ConstructCmdMessenger{OutputChannel: g}
	})

	Describe("ArgumentsNotProvided", func() {
		It("should output an appropriate error", func() {
			cm.ArgumentsNotProvided()
			Eventually(g).Should(Say("Not all required parameters were provided. See stembuild --help for more details"))
		})
	})

	Describe("LGPONotFound", func() {
		It("should output an appropriate error", func() {
			cm.LGPONotFound()
			Eventually(g).Should(Say("Could not find LGPO.zip in the current directory"))
		})
	})

	Describe("CannotConnectToVM", func() {
		It("should output an appropriate error", func() {
			connectionError := errors.New("some connection error")
			cm.CannotConnectToVM(connectionError)
			Eventually(g).Should(Say("Cannot connect to VM: %s", connectionError))
		})
	})
})
