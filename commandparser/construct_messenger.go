package commandparser

import (
	"fmt"
	"io"
)

type ConstructCmdMessenger struct {
	OutputChannel io.Writer
}

func (m *ConstructCmdMessenger) printMessage(message string) {
	fmt.Fprint(m.OutputChannel, message)
}

func (m *ConstructCmdMessenger) ArgumentsNotProvided() {
	m.printMessage("Not all required parameters were provided. See stembuild --help for more details")
}

func (m *ConstructCmdMessenger) InvalidStemcellVersion() {
	m.printMessage("Invalid stemcell version provided")
}

func (m *ConstructCmdMessenger) LGPONotFound() {
	m.printMessage("Could not find LGPO.zip in the current directory")
}
