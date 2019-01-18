package iaas_cli

import (
	_ "github.com/vmware/govmomi/govc/about"
	"github.com/vmware/govmomi/govc/cli"
	_ "github.com/vmware/govmomi/govc/object"
)

//go:generate counterfeiter . CliRunner
type CliRunner interface {
	Run(args []string) int
}

type GovcRunner struct {
}

func (r GovcRunner) Run(args []string) int {
	return cli.Run(args)
}
