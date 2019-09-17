package commandparser

import (
	"fmt"
	"io"
)

type PackageMessenger struct {
	Output io.Writer
}

func (m *PackageMessenger) InvalidOutputConfig(e error) {
	fmt.Fprintln(m.Output, e)
}

func (m *PackageMessenger) CannotCreatePackager(e error) {
	fmt.Fprintln(m.Output, e)
}

func (m *PackageMessenger) DoesNotHaveEnoughSpace(e error) {
	fmt.Fprintln(m.Output, e)
}

func (m *PackageMessenger) SourceParametersAreInvalid(e error) {
	fmt.Fprintln(m.Output, e)
}

func (m *PackageMessenger) PackageFailed(e error) {
	fmt.Fprintln(m.Output, e)
	fmt.Fprintln(m.Output, "Please provide the error logs to bosh-windows-eng@pivotal.io")
}
