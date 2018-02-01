package stemcell_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// . "github.com/onsi/gomega/gbytes"
	// . "github.com/onsi/gomega/gexec"

	"github.com/pivotal-cf-experimental/stembuild/stembuildoptions"
	"github.com/pivotal-cf-experimental/stembuild/stemcell"
)

var _ = Describe("Stemcell", func() {
	Describe("ApplyPatch", func() {

		var tmpDir string

		BeforeEach(func() {
			var err error
			tmpDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(os.RemoveAll(tmpDir)).To(Succeed())
		})
		It("successfully applies a patch", func() {
			stembuildConfig := stembuildoptions.StembuildOptions{
				PatchFile: filepath.Join("..", "testdata", "diff.patch"),
				OSVersion: "2012R2",
				Version:   "1200.1",
				VHDFile:   filepath.Join("..", "testdata", "original.vhd"),
			}
			c := stemcell.Config{
				Stop:         make(chan struct{}),
				Debugf:       func(format string, a ...interface{}) {},
				BuildOptions: stembuildConfig,
			}

			actualVmdkFilepath := filepath.Join(tmpDir, "image-disk1.vmdk")
			err := c.ApplyPatch(c.BuildOptions.VHDFile, c.BuildOptions.PatchFile, actualVmdkFilepath)
			Expect(err).NotTo(HaveOccurred())

			actualVmdk, err := ioutil.ReadFile(actualVmdkFilepath)
			fmt.Fprintf(GinkgoWriter, "image disk1: %s", actualVmdkFilepath)
			Expect(err).NotTo(HaveOccurred())

			expectedVmdkFilepath := filepath.Join("..", "testdata", "expected.vmdk")
			expectedVmdk, err := ioutil.ReadFile(expectedVmdkFilepath)
			Expect(err).NotTo(HaveOccurred())

			Expect(bytes.Equal(actualVmdk, expectedVmdk)).To(BeTrue())
		})
	})
})
