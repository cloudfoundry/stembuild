// +build !windows

package integration_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/stembuild/test/helpers"
)

var _ = Describe("Interrupts", func() {
	Describe("catchInterruptSignal", func() {
		It("cleans up on one interrupt", func() {
			var err error
			stembuildExecutable, err = helpers.BuildStembuild("1200.0.0")
			Expect(err).ToNot(HaveOccurred())

			inputVmdk := filepath.Join("..", "test", "data", "expected.vmdk")
			tmpDir, err := ioutil.TempDir(os.TempDir(), "stembuild-interrupts")
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "package", "--vmdk", inputVmdk, "--outputDir", tmpDir)
			time.Sleep(1 * time.Second)

			err = session.Command.Process.Signal(os.Interrupt)
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(1 * time.Second)

			stdErr := session.Err.Contents()
			Expect(string(stdErr)).To(ContainSubstring("received ("))
		})

		// Tried to create test to handle 2 interrupts in a row, but timing of processes makes it difficult
		// to test
	})
})
