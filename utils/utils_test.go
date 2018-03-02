package utils_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pivotal-cf-experimental/stembuild/helpers"
	. "github.com/pivotal-cf-experimental/stembuild/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

const testFile = "testdata/valid-apply-patch.yml"

var _ = Describe("Utils", func() {
	Describe("DownloadFileFromURL", func() {
		var (
			downloadPath           string
			downloadedFileFullPath string
			urlPath                string
			server                 *Server

			executeErr error
		)

		JustBeforeEach(func() {
			downloadedFileFullPath, executeErr = DownloadFileFromURL(downloadPath, urlPath, func(string, ...interface{}) {})
		})

		BeforeEach(func() {
			downloadPath = ""
			urlPath = ""
		})

		Context("when provided a valid URL", func() {
			BeforeEach(func() {
				var err error

				downloadPath, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())

				server, urlPath = helpers.StartFileServer(testFile)
			})

			AfterEach(func() {
				Expect(os.RemoveAll(downloadPath)).To(Succeed())
				server.Close()
			})

			Context("when provided a valid download path", func() {
				It("downloads the file from the URL and returns its full local path", func() {
					Expect(executeErr).NotTo(HaveOccurred())
					Expect(downloadedFileFullPath).NotTo(BeEmpty())
					Expect(helpers.CompareFiles(testFile, downloadedFileFullPath)).To(BeTrue())
				})
			})

			Context("when provided an invalid download path", func() {
				BeforeEach(func() {
					downloadPath = filepath.Join("foo", "bar")
				})

				It("fails with an appropriate error", func() {
					Expect(executeErr).To(MatchError(fmt.Sprintf("Could not create create downloaded file in directory %s", downloadPath)))
					Expect(downloadedFileFullPath).To(BeEmpty())
				})
			})
		})

		Context("when provided an invalid URL", func() {
			Context("and the URL is malformed", func() {
				BeforeEach(func() {
					var err error

					downloadPath, err = ioutil.TempDir("", "")
					Expect(err).NotTo(HaveOccurred())

					urlPath = "h:t:t://ps/foo"
				})

				It("fails with an appropriate error", func() {
					Expect(executeErr).To(MatchError(`Get h:t:t://ps/foo: unsupported protocol scheme "h"`))
					Expect(downloadedFileFullPath).To(BeEmpty())
				})
			})

			Context("but the URL is formed correctly", func() {
				BeforeEach(func() {
					var err error

					downloadPath, err = ioutil.TempDir("", "")
					Expect(err).NotTo(HaveOccurred())

					server, urlPath = helpers.StartInvalidFileServer(http.StatusNotFound)
				})

				It("fails with an appropriate error", func() {
					Expect(executeErr).To(MatchError(fmt.Sprintf("Could not create stemcell from %s\\nUnexpected response code: %d", urlPath, http.StatusNotFound)))
					Expect(downloadedFileFullPath).To(BeEmpty())
				})
			})
		})
	})

	DescribeTable("ValidateVersion",
		func(version string, isErrorExpected bool) {
			Expect(ValidateVersion(version) == nil).To(Equal(isErrorExpected))
		},
		Entry("failed to validate version", "1.2", true),
		Entry("expected error for version", "-1.2", false),
		Entry("expected error for version", "1.-2", false),
		Entry("failed to validate version", "001.002", true),
		Entry("expected error for version", "0a1.002", false),
		Entry("expected error for version", "1.a", false),
		Entry("expected error for version", "a1.2", false),
		Entry("expected error for version", "a.2", false),
		Entry("expected error for version", "1.2 a", false),
		Entry("failed to validate version", "1200.0.3-build.2", true),
		Entry("expected error for version", "1200.0.3-build.a", false),
		Entry("failed to validate version", "1.2-build.1", true),
	)
})
