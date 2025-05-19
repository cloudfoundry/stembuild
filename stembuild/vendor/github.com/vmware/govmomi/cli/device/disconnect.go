// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package device

import (
	"context"
	"flag"
	"fmt"

	"github.com/vmware/govmomi/cli"
	"github.com/vmware/govmomi/cli/flags"
)

type disconnect struct {
	*flags.VirtualMachineFlag
}

func init() {
	cli.Register("device.disconnect", &disconnect{})
}

func (cmd *disconnect) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.VirtualMachineFlag, ctx = flags.NewVirtualMachineFlag(ctx)
	cmd.VirtualMachineFlag.Register(ctx, f)
}

func (cmd *disconnect) Process(ctx context.Context) error {
	if err := cmd.VirtualMachineFlag.Process(ctx); err != nil {
		return err
	}
	return nil
}

func (cmd *disconnect) Usage() string {
	return "DEVICE..."
}

func (cmd *disconnect) Description() string {
	return `Disconnect DEVICE on VM.

Examples:
  govc device.disconnect -vm $name cdrom-3000`
}

func (cmd *disconnect) Run(ctx context.Context, f *flag.FlagSet) error {
	vm, err := cmd.VirtualMachine()
	if err != nil {
		return err
	}

	if vm == nil {
		return flag.ErrHelp
	}

	devices, err := vm.Device(ctx)
	if err != nil {
		return err
	}

	for _, name := range f.Args() {
		device := devices.Find(name)
		if device == nil {
			return fmt.Errorf("device '%s' not found", name)
		}

		if err = devices.Disconnect(device); err != nil {
			return err
		}

		if err = vm.EditDevice(ctx, device); err != nil {
			return err
		}
	}

	return nil
}
