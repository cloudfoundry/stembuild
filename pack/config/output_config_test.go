package config_test

import (
	"os"
	"path/filepath"

	. "github.com/cloudfoundry-incubator/stembuild/pack/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OutputConfig", func() {
	Describe("os", func() {
		Context("no os specified", func() {
			It("should be invalid", func() {
				valid := IsValidOS("")
				Expect(valid).To(BeFalse())
			})
		})
		Context("a supported os is specified", func() {
			It("2016 should be valid", func() {
				valid := IsValidOS("2016")
				Expect(valid).To(BeTrue())
			})
			It("2012R2 should be valid", func() {
				valid := IsValidOS("2012R2")
				Expect(valid).To(BeTrue())
			})
			It("1803 should be valid", func() {
				valid := IsValidOS("1803")
				Expect(valid).To(BeTrue())
			})
		})
		Context("something other than a supported os is specified", func() {
			It("should be invalid", func() {
				valid := IsValidOS("bad-thing")
				Expect(valid).To(BeFalse())
			})
			It("1709 should be invalid", func() {
				valid := IsValidOS("1709")
				Expect(valid).To(BeFalse())
			})
		})

	})
	Describe("stemcell version", func() {
		Context("no stemcell version specified", func() {
			It("should be invalid", func() {
				valid := IsValidStemcellVersion("")
				Expect(valid).To(BeFalse())
			})
		})
		Context("stemcell version specified is valid pattern", func() {
			It("should be valid", func() {
				valids := []string{"4.4", "4.4-build.1", "4.4.4", "4.4.4-build.4"}
				for _, version := range valids {
					valid := IsValidStemcellVersion(version)
					Expect(valid).To(BeTrue())
				}
			})
		})
		Context("stemcell version specified is invalid pattern", func() {
			It("should be invalid", func() {
				valid := IsValidStemcellVersion("4.4.4.4")
				Expect(valid).To(BeFalse())
			})
		})
	})
	Describe("validateOutputDir", func() {

		Context("no dir given", func() {
			It("should be an error", func() {
				err := ValidateOrCreateOutputDir("")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("valid relative directory that does not exist", func() {
			It("should create the directory and not fail", func() {
				err := os.RemoveAll(filepath.Join(".", "tmp"))
				Expect(err).ToNot(HaveOccurred())
				err = ValidateOrCreateOutputDir(filepath.Join(".", "tmp"))
				Expect(err).ToNot(HaveOccurred())
				_, err = os.Stat(filepath.Join(filepath.Join(".", "tmp")))
				Expect(err).ToNot(HaveOccurred())
			})
			AfterEach(func() {
				err := os.RemoveAll(filepath.Join(".", "tmp"))
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("valid directory that already exists", func() {
			It("should not fail", func() {
				err := os.Mkdir(filepath.Join(".", "tmp"), 0700)
				Expect(err).ToNot(HaveOccurred())
				err = ValidateOrCreateOutputDir(filepath.Join(".", "tmp"))
				Expect(err).ToNot(HaveOccurred())
			})
			AfterEach(func() {
				err := os.RemoveAll(filepath.Join(".", "tmp"))
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("valid absolute directory", func() {
			It("should not fail", func() {
				cwd, err := os.Getwd()
				Expect(err).ToNot(HaveOccurred())
				err = ValidateOrCreateOutputDir(cwd)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("invalid directory because intermediate directories do not exist", func() {
			It("should be an error", func() {
				err := os.RemoveAll(filepath.Join(".", "tmp"))
				Expect(err).ToNot(HaveOccurred())
				err = ValidateOrCreateOutputDir(filepath.Join(".", "tmp", "subtmp"))
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
