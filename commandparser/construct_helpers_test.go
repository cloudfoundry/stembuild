package commandparser_test

import (
	. "github.com/cloudfoundry-incubator/stembuild/commandparser"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var _ = Describe("construct_helpers", func() {

	Describe("IsArtifactInDirectory", func() {
		Context("Directory given is valid", func() {
			Describe("automation artifact", func() {
				filename := "StemcellAutomation.zip"

				Context("artifact is not present", func() {
					dir := filepath.Join("..", "test", "constructData", "emptyDir")

					It("should return false with no error", func() {
						present, err := IsArtifactInDirectory(dir, filename)
						Expect(err).ToNot(HaveOccurred())
						Expect(present).To(BeFalse())
					})
				})

				Context("artifact is present", func() {
					dir := filepath.Join("..", "test", "constructData", "fullDir")

					It("should return true with no error", func() {
						present, err := IsArtifactInDirectory(dir, filename)
						Expect(err).ToNot(HaveOccurred())
						Expect(present).To(BeTrue())
					})
				})
			})

			Describe("LGPO", func() {
				filename := "LGPO.zip"

				Context("LGPO is not present", func() {
					dir := filepath.Join("..", "test", "constructData", "emptyDir")

					It("should return false with no error", func() {
						present, err := IsArtifactInDirectory(dir, filename)
						Expect(err).ToNot(HaveOccurred())
						Expect(present).To(BeFalse())
					})
				})

				Context("artifact is present", func() {
					dir := filepath.Join("..", "test", "constructData", "fullDir")

					It("should return true with no error", func() {
						present, err := IsArtifactInDirectory(dir, filename)
						Expect(err).ToNot(HaveOccurred())
						Expect(present).To(BeTrue())
					})
				})
			})
		})

		Context("Directory given is not valid", func() {
			filename := "file"
			It("should return an error", func() {
				dir := filepath.Join("..", "test", "constructData", "notExist")
				_, err := IsArtifactInDirectory(dir, filename)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
