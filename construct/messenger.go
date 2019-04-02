package construct

import "io"

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
