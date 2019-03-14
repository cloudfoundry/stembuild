package vmconstruct_factory

import (
	"github.com/cloudfoundry-incubator/stembuild/commandparser"
	"github.com/cloudfoundry-incubator/stembuild/construct"
)

type VMConstructFactory struct {
}

func (f *VMConstructFactory) GetVMPreparer(winrmIp string, winrmUsername string, winrmPassword string) commandparser.VMPreparer {
	return construct.NewVMConstruct(winrmIp, winrmUsername, winrmPassword)
}
