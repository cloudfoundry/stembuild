package remotemanager

import (
	"errors"
	"fmt"
	"github.com/cloudfoundry-incubator/stembuild/poller"
	"time"
)

var tryCheckReboot = `shutdown /r /f /t 60 /c "packer restart test"`
var abortReboot = `shutdown /a`

type RebootWaiter struct {
	poller        poller.PollerI
	rebootChecker RebootCheckerI
}

func NewRebootWaiter(poller poller.PollerI, rebootChecker RebootCheckerI) *RebootWaiter {
	return &RebootWaiter{
		poller,
		rebootChecker,
	}
}

func (rw *RebootWaiter) WaitForRebootFinished() error {
	err := rw.poller.Poll(10*time.Second, rw.rebootChecker.RebootHasFinished)

	if err != nil {
		return fmt.Errorf("error polling for reboot: %s", err)
	}
	return nil
}

//go:generate counterfeiter . RebootCheckerI
type RebootCheckerI interface {
	RebootHasFinished() (bool, error)
}

type RebootChecker struct {
	remoteManager RemoteManager
}

func NewRebootChecker(winrmRemoteManager RemoteManager) *RebootChecker {
	return &RebootChecker{winrmRemoteManager}
}

func (rc *RebootChecker) RebootHasFinished() (bool, error) {

	exitCode, err := rc.remoteManager.ExecuteCommand(tryCheckReboot)
	if err != nil {
		return false, nil
	}
	if exitCode == 0 {
		var abortExitCode int
		var abortErr error
		for i := 0; i < 5; i++ {
			abortExitCode, abortErr = rc.remoteManager.ExecuteCommand(abortReboot)

			if abortErr == nil{
				break
			}
		}

		if abortErr != nil {
			return false, fmt.Errorf("unable to abort reboot: %s", abortErr)
		}

		if abortExitCode == 0 {
			return true, nil
		} else  {
			return false, errors.New("unable to abort reboot.")
		}

	}
	return false, err
}
