// +build !darwin,!windows

package ovftool_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/stembuild/ovftool"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Describe("ovftool", func() {

	var oldPath string

	BeforeEach(func() {
		oldPath = os.Getenv("PATH")
		os.Unsetenv("PATH")
	})

	AfterEach(func() {
		os.Setenv("PATH", oldPath)
	})

	It("when ovftool is on the PATH, its path is returned", func() {
		tmpDir, err := ioutil.TempDir(os.TempDir(), "ovftmp")
		Expect(err).ToNot(HaveOccurred())
		err = ioutil.WriteFile(filepath.Join(tmpDir, "ovftool"), []byte{}, 0700)
		Expect(err).ToNot(HaveOccurred())
		os.Setenv("PATH", tmpDir)

		ovfPath, err := ovftool.Ovftool([]string{})
		os.RemoveAll(tmpDir)

		Expect(err).ToNot(HaveOccurred())
		Expect(ovfPath).To(Equal(filepath.Join(tmpDir, "ovftool")))
	})

	It("fails when ovftool is not installed in the PATH", func() {
		_, err := ovftool.Ovftool([]string{})
		Expect(err).To(HaveOccurred())
	})
})
