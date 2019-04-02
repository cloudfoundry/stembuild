package construct

import (
	"fmt"
	"io"
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

func (m *Messenger) ExecuteScriptStarted() {
	m.out.Write([]byte("\nExecuting setup script...\n"))
}

func (m *Messenger) ExecuteScriptSucceeded() {
	m.out.Write([]byte("\nFinished executing setup script.\n"))
}

func (m *Messenger) UploadFileStarted(artifact string) {
	m.out.Write([]byte(fmt.Sprintf("\tUploading %s to target VM...", artifact)))
}

func (m *Messenger) UploadFileSucceeded() {
	m.out.Write([]byte("succeeded.\n"))
}
