package stembuild_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/stembuild"
	"io/ioutil"
	"os"
)

var _ = Describe("inputs", func() {
	Describe("vmdk", func() {
		Context("no vmdk specified", func() {
			vmdk := ""
			It("should be invalid", func() {
				valid, err := stembuild.IsValidVMDK(vmdk)
				Expect(err).ToNot(HaveOccurred())
				Expect(valid).To(BeFalse())
			})
		})
		Context("valid vmdk file specified", func() {
			It("should be valid", func() {

				vmdk, err := ioutil.TempFile("", "temp.vmdk")
				Expect(err).ToNot(HaveOccurred())
				defer os.Remove(vmdk.Name())

				valid, err := stembuild.IsValidVMDK(vmdk.Name())
				Expect(err).To(BeNil())
				Expect(valid).To(BeTrue())
			})
		})
		Context("invalid vmdk file specified", func() {
			It("should be invalid", func() {
				valid, err := stembuild.IsValidVMDK("/dev/null")
				Expect(err).To(BeNil())
				Expect(valid).To(BeFalse())
			})
		})
	})
	Describe("os", func() {
		Context("no os specified", func() {
			It("should be invalid", func() {
				valid := stembuild.IsValidOS("")
				Expect(valid).To(BeFalse())
			})
		})
		Context("a supported os is specified", func() {
			It("should be valid", func() {
				valid := stembuild.IsValidOS("1709")
				Expect(valid).To(BeTrue())
			})
		})
		Context("something other than a supported os is specified", func() {
			It("should be invalid", func() {
				valid := stembuild.IsValidOS("bad-thing")
				Expect(valid).To(BeFalse())
			})
		})

	})
	Describe("version", func() {
		Context("no version specified", func() {
			It("should be invalid", func() {
				valid := stembuild.IsValidVersion("")
				Expect(valid).To(BeFalse())
			})
		})
		Context("version specified is valid pattern", func() {
			It("should be valid", func() {
				valids := []string{"4.4", "4.4-build.1", "4.4.4", "4.4.4-build.4"}
				for _, version := range valids {
					valid := stembuild.IsValidVersion(version)
					Expect(valid).To(BeTrue())
				}
			})
		})
		Context("version specified is invalid pattern", func() {
			It("should be invalid", func() {
				valid := stembuild.IsValidVersion("4.4.4.4")
				Expect(valid).To(BeFalse())
			})
		})
	})
})
