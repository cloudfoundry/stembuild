package ovftool_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

func vmwareInstallPaths() ([]string, error) {
	return ["path/to/ovftool"], nil
}

func lookPath(name string) (string, err) {
	return "", err.Error("Hello");
}

func findExecutable(dir, name string) (string, err) {
}

var _ = Describe("ovftool", func() {

	Context("no ovftool found", func() {
		BeforeEach(func() {
			
		})

		path, err := ovftool.Ovftool()
		Expect(path).To(Equal(""))
		Expect(err).To(HaveOccurred())
	})
	Context("ovftool found", func() {

	})
})
