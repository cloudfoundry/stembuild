package construct_test

import (
	. "github.com/cloudfoundry-incubator/stembuild/remotemanager"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var _ = Describe("WinRM Remote Manager", func() {

	var rm RemoteManager

	BeforeEach(func() {
		rm = NewWinRM(conf.TargetIP, conf.VMUsername, conf.VMPassword)
		Expect(rm).ToNot(BeNil())
	})

	Context("ExtractArchive", func() {
		BeforeEach(func() {
			err := rm.UploadArtifact(filepath.Join("assets", "StemcellAutomation.zip"), "C:\\provision\\StemcellAutomation.zip")
			Expect(err).ToNot(HaveOccurred())

		})

		It("succeeds when Extract-Archive powershell function returns zero exit code", func() {
			err := rm.ExtractArchive("C:\\provision\\StemcellAutomation.zip", "C:\\provision")
			Expect(err).ToNot(HaveOccurred())
		})

		It("fails when Extract-Acrhive powershell function returns non-zero exit code", func() {
			err := rm.ExtractArchive("C:\\provision\\NonExistingFile.zip", "C:\\provision")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("powershell encountered an issue: "))
		})

		AfterEach(func() {
			err := rm.ExecuteCommand("powershell.exe Remove-Item c:\\provision -recurse")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("ExecuteCommand", func() {

		It("succeeds when powershell command returns a zero exit code", func() {
			err := rm.ExecuteCommand("powershell.exe ls c:\\windows")
			Expect(err).ToNot(HaveOccurred())
		})

		It("fails when powershell command returns non-zero exit code", func() {
			err := rm.ExecuteCommand("powershell.exe notRealCommand")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("powershell encountered an issue: "))
		})

	})
})
