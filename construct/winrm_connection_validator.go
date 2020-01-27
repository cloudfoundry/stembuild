package construct

import (
	. "github.com/cloudfoundry-incubator/stembuild/remotemanager"
)

type WinRMConnectionValidator struct {
	RemoteManager RemoteManager
}

func (v *WinRMConnectionValidator) Validate() error {
	err := v.RemoteManager.CanReachVM()
	if err != nil {
		return err
	}

	err = v.RemoteManager.CanLoginVM()
	if err != nil {
		return err
	}

	return nil
}
