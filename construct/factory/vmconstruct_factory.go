package vmconstruct_factory

import (
	"context"
	"os"

	p "github.com/cloudfoundry/stembuild/poller"

	"github.com/cloudfoundry/stembuild/version"

	"github.com/cloudfoundry/stembuild/commandparser"
	"github.com/cloudfoundry/stembuild/construct"
	"github.com/cloudfoundry/stembuild/construct/archive"
	"github.com/cloudfoundry/stembuild/construct/config"
	"github.com/cloudfoundry/stembuild/iaas_cli"
	"github.com/cloudfoundry/stembuild/iaas_cli/iaas_clients"
	"github.com/pkg/errors"

	. "github.com/cloudfoundry/stembuild/remotemanager"
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
		return nil, errors.Wrap(err, "Cannot complete login due to an incorrect vCenter user name or password")
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

	winRmClientFactory := NewWinRmClientFactory(config.GuestVmIp, config.GuestVMUsername, config.GuestVMPassword)
	remoteManager := NewWinRM(config.GuestVmIp, config.GuestVMUsername, config.GuestVMPassword, winRmClientFactory)

	vmConnectionValidator := &construct.WinRMConnectionValidator{
		RemoteManager: remoteManager,
	}

	poller := &p.Poller{}

	rebootChecker := NewRebootChecker(remoteManager)

	rebootWaiter := NewRebootWaiter(poller, rebootChecker)

	scriptExecutor := construct.NewScriptExecutor(remoteManager)

	return construct.NewVMConstruct(
		ctx,
		remoteManager,
		config.GuestVMUsername,
		config.GuestVMPassword,
		config.VmInventoryPath,
		client,
		guestManager,
		winRMManager,
		vmConnectionValidator,
		messenger,
		poller,
		versionGetter,
		rebootWaiter,
		scriptExecutor,
		config.SetupFlags,
	), nil
}
