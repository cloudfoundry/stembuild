package construct

import (
	"fmt"
	"io"
	"time"
)

type Messenger struct {
	out io.Writer
}

func NewMessenger(out io.Writer) *Messenger {
	return &Messenger{out}
}

func (m *Messenger) EnableWinRMStarted() {
	m.out.Write([]byte("\nAttempting to enable WinRM on the guest vm..."))
}

func (m *Messenger) EnableWinRMSucceeded() {
	m.out.Write([]byte("WinRm enabled on the guest VM\n"))
}

func (m *Messenger) ValidateVMConnectionStarted() {
	m.out.Write([]byte("\nValidating connection to vm..."))
}

func (m *Messenger) ValidateVMConnectionSucceeded() {
	m.out.Write([]byte("succeeded.\n"))
}

func (m *Messenger) CreateProvisionDirStarted() {
	m.out.Write([]byte("\nCreating provision dir on target VM..."))
}

func (m *Messenger) CreateProvisionDirSucceeded() {
	m.out.Write([]byte("succeeded.\n"))
}

func (m *Messenger) UploadArtifactsStarted() {
	m.out.Write([]byte("\nTransferring ~20 MB to the Windows VM. Depending on your connection, the transfer may take 15-45 minutes\n"))
}

func (m *Messenger) UploadArtifactsSucceeded() {
	m.out.Write([]byte("\nAll files have been uploaded.\n"))
}

func (m *Messenger) ExtractArtifactsStarted() {
	m.out.Write([]byte("\nExtracting artifacts..."))
}

func (m *Messenger) ExtractArtifactsSucceeded() {
	m.out.Write([]byte("succeeded.\n"))
}

func (m *Messenger) ExecuteSetupScriptStarted() {
	m.out.Write([]byte("\nExecuting setup script 1 of 2...\n"))
}

func (m *Messenger) ExecuteSetupScriptSucceeded() {
	m.out.Write([]byte("\nFinished executing setup script 1 of 2.\n"))
}

func (m *Messenger) RebootHasStarted() {
	m.out.Write([]byte("\nThe reboot has started...\n"))
}

func (m *Messenger) RebootHasFinished() {
	m.out.Write([]byte("\nThe reboot has finished.\n"))
}

func (m *Messenger) ExecutePostRebootScriptStarted() {
	m.out.Write([]byte("\nExecuting setup script 2 of 2...\n"))
}

func (m *Messenger) ExecutePostRebootScriptSucceeded() {
	m.out.Write([]byte("\nFinished executing setup script 2 of 2.\n"))

}

func (m *Messenger) ExecutePostRebootWarning(warning string) {
	m.out.Write([]byte("\n"))
	m.out.Write([]byte(warning))
	m.out.Write([]byte("\n"))
}

func (m *Messenger) UploadFileStarted(artifact string) {
	m.out.Write([]byte(fmt.Sprintf("\tUploading %s to target VM...", artifact)))
}

func (m *Messenger) UploadFileSucceeded() {
	m.out.Write([]byte("succeeded.\n"))
}

func (m *Messenger) OSVersionFileCreationFailed(errorMessage string) {
	m.logValidateOSWarning("OS Version file creation failed", errorMessage)
}

func (m *Messenger) ExitCodeRetrievalFailed(errorMessage string) {
	m.logValidateOSWarning("Failed to retrieve exit code for process to create OS Version file", errorMessage)
}

func (m *Messenger) DownloadFileFailed(errorMessage string) {
	m.logValidateOSWarning("Failed to download OS Version file", errorMessage)
}

func (m *Messenger) logValidateOSWarning(log string, errorMessage string) {
	matchingVersionWarning := "Ensure the version of the stemcell you're trying to build matches the corresponding base ISO you're using.\n" +
		"For example: If you're building 2019.x, then you should be using 'Windows Server 2019' only"
	m.out.Write([]byte(fmt.Sprintf("Warning: %s:\n%s\n%s", log, matchingVersionWarning, errorMessage)))
}

func (m *Messenger) WaitingForShutdown() {
	t := time.Now()
	timeStampFormat := "2006-01-02T15:04:05.999999-07:00"
	m.out.Write([]byte(fmt.Sprintf("%s Still preparing VM...\n", t.Format(timeStampFormat))))
}

func (m *Messenger) ShutdownCompleted() {
	m.out.Write([]byte("VM has now been shutdown. Run `stembuild package` to finish building the stemcell.\n"))
}

func (m *Messenger) WinRMDisconnectedForReboot() {
	m.out.Write([]byte("\nWinRM has been disconnected so the VM can reboot.\n"))

}
