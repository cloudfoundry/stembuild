package iaas_cli

import (
	"bytes"
	"io"
	"os"

	_ "github.com/vmware/govmomi/govc/about"
	"github.com/vmware/govmomi/govc/cli"
	_ "github.com/vmware/govmomi/govc/device"
	_ "github.com/vmware/govmomi/govc/device/cdrom"
	_ "github.com/vmware/govmomi/govc/export"
	_ "github.com/vmware/govmomi/govc/object"
)

//go:generate counterfeiter . CliRunner
type CliRunner interface {
	Run(args []string) int
	RunWithOutput(args []string) (string, int, error)
}

type GovcRunner struct {
}

func (r *GovcRunner) Run(args []string) int {
	return cli.Run(args)
}

func (r *GovcRunner) RunWithOutput(args []string) (string, int, error) {
	old := os.Stdout // keep backup of the real stdout
	reader, w, _ := os.Pipe()
	os.Stdout = w

	print()

	outC := make(chan string)
	errC := make(chan error)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, reader)
		if err != nil {
			errC <- err
		} else {
			outC <- buf.String()
		}
	}()

	exitCode := r.Run(args)

	// back to normal state
	err := w.Close()
	os.Stdout = old // restoring the real stdout
	if err != nil {
		return "", exitCode, err
	}

	select {
	case out := <-outC:
		return out, exitCode, nil
	case err := <-errC:
		return "", exitCode, err
	}
}
