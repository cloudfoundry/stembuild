package vmconstruct_factory

import (
	"os"

	"github.com/cloudfoundry-incubator/stembuild/commandparser"
	"github.com/cloudfoundry-incubator/stembuild/construct"
	"github.com/cloudfoundry-incubator/stembuild/construct/archive"
	"github.com/cloudfoundry-incubator/stembuild/construct/config"
	"github.com/cloudfoundry-incubator/stembuild/iaas_cli"
	"github.com/cloudfoundry-incubator/stembuild/iaas_cli/iaas_clients"
)

type VMConstructFactory struct {
}

func (f *VMConstructFactory) VMPreparer(config config.SourceConfig) commandparser.VmConstruct {
	runner := &iaas_cli.GovcRunner{}
	client := iaas_clients.NewVcenterClient(config.VCenterUsername, config.VCenterPassword, config.VCenterUrl, config.CaCertFile, runner)

	zip := &archive.Zip{}

	messenger := construct.NewMessenger(os.Stdout)

	return construct.NewVMConstruct(
		config.GuestVmIp,
		config.GuestVMUsername,
		config.GuestVMPassword,
		config.VmInventoryPath,
		client,
		zip,
		messenger,
	)
}
