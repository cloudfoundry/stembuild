// +build darwin

package ovftool_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/ovftool"
)

var _ = Describe("ovftool darwin", func() {

	Context("homeDirectory", func() {
		var home = ""

		BeforeEach(func() {
			home = os.Getenv("HOME")
		})

		AfterEach(func() {
			os.Setenv("HOME", home)
		})

		It("returns value of HOME environment variable is set", func() {
			Expect(home).NotTo(Equal(""))

			testHome := ovftool.HomeDirectory()
			Expect(testHome).To(Equal(home))
		})

		It("returns user HOME directory if HOME environment variable is not set", func() {
			os.Unsetenv("HOME")

			testHome := ovftool.HomeDirectory()
			Expect(testHome).To(Equal(home))
		})

	})

	Context("Ovftool", func() {
		var oldPath string

		BeforeEach(func() {
			oldPath = os.Getenv("PATH")
			os.Unsetenv("PATH")
		})

		AfterEach(func() {
			os.Setenv("PATH", oldPath)
		})

		It("when ovftool is on the path, its path is returned", func() {
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

		It("fails when given invalid set of install paths", func() {
			tmpDir, err := ioutil.TempDir(os.TempDir(), "ovftmp")
			Expect(err).ToNot(HaveOccurred())

			_, err = ovftool.Ovftool([]string{tmpDir})
			os.RemoveAll(tmpDir)

			Expect(err).To(HaveOccurred())
		})

		It("fails when given empty set of install paths", func() {
			_, err := ovftool.Ovftool([]string{})
			Expect(err).To(HaveOccurred())
		})

		Context("when ovftool is installed", func() {
			var tmpDir, dummyDir string

			BeforeEach(func() {
				var err error
				tmpDir, err = ioutil.TempDir(os.TempDir(), "ovftmp")
				Expect(err).ToNot(HaveOccurred())
				dummyDir, err = ioutil.TempDir(os.TempDir(), "trashdir")
				Expect(err).ToNot(HaveOccurred())
				err = ioutil.WriteFile(filepath.Join(tmpDir, "ovftool"), []byte{}, 0700)
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				os.RemoveAll(tmpDir)
				os.RemoveAll(dummyDir)
			})

			It("Returns the path to the ovftool", func() {
				ovfPath, err := ovftool.Ovftool([]string{"notrealdir", dummyDir, tmpDir})

				Expect(err).ToNot(HaveOccurred())
				Expect(ovfPath).To(Equal(filepath.Join(tmpDir, "ovftool")))
			})
		})
	})
})
