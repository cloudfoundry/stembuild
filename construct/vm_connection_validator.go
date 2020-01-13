package construct

import (
	. "github.com/cloudfoundry-incubator/stembuild/remotemanager"
	"github.com/pkg/errors"
)

type VMConnectionValidator struct {
	RemoteManager RemoteManager
}

func (v *VMConnectionValidator) Validate() error {
	err := v.RemoteManager.CanReachVM()
	if err != nil {
		return err
	}

	err = v.RemoteManager.CanLoginVM()
	if err != nil {
		return errors.Wrap(err, "Cannot complete login due to an incorrect VM user name or password")
	}

	return nil
}
