package guest_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/vmware/govmomi/vim25/types"
)

//go:generate counterfeiter . ProcManager
type ProcManager interface {
	StartProgram(ctx context.Context, auth types.BaseGuestAuthentication, spec types.BaseGuestProgramSpec) (int64, error)
	ListProcesses(ctx context.Context, auth types.BaseGuestAuthentication, pids []int64) ([]types.GuestProcessInfo, error)
}

type GuestManager struct {
	auth           types.NamePasswordAuthentication
	processManager ProcManager
}

func NewGuestManager(auth types.NamePasswordAuthentication, processManager ProcManager) *GuestManager {
	return &GuestManager{auth, processManager}
}

func (g *GuestManager) StartProgramInGuest(ctx context.Context, command, args string) (int64, error) {
	spec := types.GuestProgramSpec{
		ProgramPath: command,
		Arguments:   args,
	}

	pid, err := g.processManager.StartProgram(ctx, &g.auth, &spec)
	if err != nil {
		return -1, fmt.Errorf("vcenter_client - could not run process: %s on guest os, error: %s",
			fmt.Sprintf("%s %s", command, args), err.Error())
	}

	return pid, nil
}

func (g *GuestManager) ExitCodeForProgramInGuest(ctx context.Context, pid int64) (int32, error) {
	for {
		procs, err := g.processManager.ListProcesses(ctx, &g.auth, []int64{pid})
		if err != nil {
			return -1, fmt.Errorf("vcenter_client - could not observe program exiting: %s", err.Error())
		}

		if len(procs) != 1 {
			return -1, fmt.Errorf("vcenter_client - could not observe program exiting")
		}

		if procs[0].EndTime == nil {
			<-time.After(time.Millisecond * 250)
			continue
		}

		return procs[0].ExitCode, nil
	}
}
