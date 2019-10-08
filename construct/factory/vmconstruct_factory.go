package vmconstruct_factory

import (
	"context"
	"os"

	"github.com/cloudfoundry-incubator/stembuild/construct/poller"

	"github.com/cloudfoundry-incubator/stembuild/version"

	"github.com/cloudfoundry-incubator/stembuild/commandparser"
	"github.com/cloudfoundry-incubator/stembuild/construct"
	"github.com/cloudfoundry-incubator/stembuild/construct/archive"
	"github.com/cloudfoundry-incubator/stembuild/construct/config"
	"github.com/cloudfoundry-incubator/stembuild/iaas_cli"
	"github.com/cloudfoundry-incubator/stembuild/iaas_cli/iaas_clients"

	. "github.com/cloudfoundry-incubator/stembuild/remotemanager"
)

type VMConstructFactory struct {
}

func (f *VMConstructFactory) VMPreparer(config config.SourceConfig, vCenterManager commandparser.VCenterManager) (commandparser.VmConstruct, error) {
	runner := &iaas_cli.GovcRunner{}
	client := iaas_clients.NewVcenterClient(config.VCenterUsername, config.VCenterPassword, config.VCenterUrl, config.CaCertFile, runner)

	messenger := construct.NewMessenger(os.Stdout)

	ctx := context.Background()
	err := vCenterManager.Login(ctx)
	if err != nil {
		return nil, err
	}

	vm, err := vCenterManager.FindVM(ctx, config.VmInventoryPath)
	if err != nil {
		return nil, err
	}

	opsManager := vCenterManager.OperationsManager(ctx, vm)

	guestManager, err := vCenterManager.GuestManager(ctx, opsManager, config.GuestVMUsername, config.GuestVMPassword)
	if err != nil {
		return nil, err
	}

	winRMManager := &construct.WinRMManager{
		GuestManager: guestManager,
		Unarchiver:   &archive.Zip{},
	}
	versionGetter := version.NewVersionGetter()
	osValidator := &construct.OSVersionValidator{
		GuestManager:  guestManager,
		Messenger:     messenger,
		VersionGetter: versionGetter,
	}

	return construct.NewVMConstruct(
		ctx,
		NewWinRM(config.GuestVmIp, config.GuestVMUsername, config.GuestVMPassword),
		config.GuestVMUsername,
		config.GuestVMPassword,
		config.VmInventoryPath,
		client,
		guestManager,
		winRMManager,
		osValidator,
		messenger,
		&poller.Poller{},
	), nil
}
