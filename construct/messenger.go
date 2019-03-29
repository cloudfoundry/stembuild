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
