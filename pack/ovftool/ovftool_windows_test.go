package ovftool_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Describe("ovftool", func() {

	Context("vmwareInstallPaths", func() {

		It("returns install paths when given valid set of registry key paths", func() {

			keypaths := []string{
				`SOFTWARE\Wow6432Node\VMware, Inc.\VMware Workstation`,
				`SOFTWARE\Wow6432Node\VMware, Inc.\VMware OVF Tool`,
				`SOFTWARE\VMware, Inc.\VMware Workstation`,
				`SOFTWARE\VMware, Inc.\VMware OVF Tool`,
			}

			searchPaths, err := ovftool.VmwareInstallPaths(keypaths)

			Expect(err).ToNot(HaveOccurred())
			Expect(searchPaths).ToNot(BeNil())
		})

		It("fails when given invalid registry key path", func() {
			keypaths := []string{`\SOFTWARE\faketempkey`}

			_, err := ovftool.VmwareInstallPaths(keypaths)

			Expect(err).To(HaveOccurred())
		})

		It("fails when given empty set of registry keys paths", func() {
			var keypaths []string

			_, err := ovftool.VmwareInstallPaths(keypaths)

			Expect(err).To(HaveOccurred())
		})

		Context("when given a valid registry keypath that does not have an installPath[64] value", func() {
			var key registry.Key
			var keypaths []string

			BeforeEach(func() {
				var err error
				key, err = registry.OpenKey(registry.CURRENT_USER, `SOFTWARE`, registry.WRITE)
				Expect(err).ToNot(HaveOccurred())
				_, exists, err := registry.CreateKey(key, `faketempkey`, registry.WRITE)
				Expect(err).ToNot(HaveOccurred())
				Expect(exists).To(BeFalse())
				keypaths = []string{`\SOFTWARE\faketempkey`}
			})

			AfterEach(func() {
				Expect(key).ToNot(BeNil())
				err := registry.DeleteKey(key, `faketempkey`)
				Expect(err).ToNot(HaveOccurred())
			})

			It("fails", func() {

				_, err := ovftool.VmwareInstallPaths(keypaths)

				Expect(err).To(HaveOccurred())
			})
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
			err = ioutil.WriteFile(filepath.Join(tmpDir, "ovftool.exe"), []byte{}, 0700)
			Expect(err).ToNot(HaveOccurred())
			os.Setenv("PATH", tmpDir)

			ovfPath, err := ovftool.Ovftool([]string{})
			os.RemoveAll(tmpDir)

			Expect(err).ToNot(HaveOccurred())
			Expect(ovfPath).To(Equal(filepath.Join(tmpDir, "ovftool.exe")))
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
				err = ioutil.WriteFile(filepath.Join(tmpDir, "ovftool.exe"), []byte{}, 0644)
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				os.RemoveAll(tmpDir)
				os.RemoveAll(dummyDir)
			})

			It("Returns the path to the ovftool", func() {
				ovfPath, err := ovftool.Ovftool([]string{"notrealdir", dummyDir, tmpDir})

				Expect(err).ToNot(HaveOccurred())
				Expect(ovfPath).To(Equal(filepath.Join(tmpDir, "ovftool.exe")))
			})
		})
	})
})
