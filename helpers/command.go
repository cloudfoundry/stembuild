package helpers

import (
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

const (
	Cmd                = "stembuild"
	DebugCommandPrefix = "\nCMD>"
	DebugOutPrefix     = "OUT: "
	DebugErrPrefix     = "ERR: "
)

func Stembuild(args ...string) *Session {
	WriteCommand(args)
	session, err := Start(
		exec.Command(Cmd, args...),
		NewPrefixedWriter(DebugOutPrefix, GinkgoWriter),
		NewPrefixedWriter(DebugErrPrefix, GinkgoWriter))
	Expect(err).NotTo(HaveOccurred())
	return session
}

func WriteCommand(args []string) {
	display := append([]string{DebugCommandPrefix, Cmd}, args...)
	GinkgoWriter.Write([]byte(strings.Join(append(display, "\n"), " ")))
}
