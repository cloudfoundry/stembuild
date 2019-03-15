package commandparser_test

import (
	"github.com/cloudfoundry-incubator/stembuild/commandparser"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("ConstructValidator", func() {

	var (
		c commandparser.ConstructValidator
	)
	BeforeEach(func() {
		c = commandparser.ConstructValidator{}
	})

	Describe("NonEmptyArgs", func() {
		It("should return true if all arguments are present", func() {
			nonEmptyArgs := c.NonEmptyArgs("version", "ip", "username", "password")
			Expect(nonEmptyArgs).To(BeTrue())
		})

		It("should return false if stemcellVersion argument is empty", func() {
			nonEmptyArgs := c.NonEmptyArgs("", "ip", "username", "password")
			Expect(nonEmptyArgs).To(BeFalse())
		})

		It("should return false if winrmIp argument is empty", func() {
			nonEmptyArgs := c.NonEmptyArgs("version", "", "username", "password")
			Expect(nonEmptyArgs).To(BeFalse())
		})

		It("should return false if username argument is empty", func() {
			nonEmptyArgs := c.NonEmptyArgs("version", "ip", "", "password")
			Expect(nonEmptyArgs).To(BeFalse())
		})

		It("should return false if password argument is empty", func() {
			nonEmptyArgs := c.NonEmptyArgs("version", "ip", "username", "")
			Expect(nonEmptyArgs).To(BeFalse())
		})
	})

	Describe("LGPOInDirectory", func() {
		It("should return true if LGPO exists in the directory", func() {
			_, err := os.Create("LGPO.zip")
			os.Stat("LGPO.zip")
			Expect(err).ToNot(HaveOccurred())

			result := c.LGPOInDirectory()

			Expect(result).To(BeTrue())

			os.Remove("LGPO.zip")
		})

		It("should return false if LGPO doesn't exist in the directory", func() {
			result := c.LGPOInDirectory()

			Expect(result).To(BeFalse())
		})
	})

	Describe("ValidStemcellInfo", func() {
		It("should return true if the given stemcell info is valid", func() {
			result := c.ValidStemcellInfo("1803.9999")

			Expect(result).To(BeTrue())
		})

		It("Should return false if the given stemcell info is invalid", func() {
			result := c.ValidStemcellInfo("completely-invalid")

			Expect(result).To(BeFalse())
		})
	})
})
