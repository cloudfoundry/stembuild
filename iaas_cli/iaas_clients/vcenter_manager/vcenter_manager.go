package vcenter_manager

import (
	"context"
	"net/url"

	"github.com/vmware/govmomi/object"

	"github.com/vmware/govmomi/vim25"

	"github.com/vmware/govmomi/guest"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/cloudfoundry-incubator/stembuild/iaas_cli/iaas_clients/guest_manager"
)

//go:generate counterfeiter . GovmomiClient
type GovmomiClient interface {
	Login(ctx context.Context, u *url.Userinfo) error
}

//go:generate counterfeiter . Finder
type Finder interface {
	VirtualMachine(ctx context.Context, path string) (*object.VirtualMachine, error)
}

//go:generate counterfeiter . OpsManager
type OpsManager interface {
	ProcessManager(ctx context.Context) (*guest.ProcessManager, error)
}

type VCenterManager struct {
	govmomiClient GovmomiClient
	vimClient     *vim25.Client
	finder        Finder
	username      string
	password      string
}

func NewVCenterManager(govmomiClient GovmomiClient, vimClient *vim25.Client, finder Finder, username, password string) (*VCenterManager, error) {
	return &VCenterManager{govmomiClient: govmomiClient, vimClient: vimClient, finder: finder, username: username, password: password}, nil
}

func (v *VCenterManager) Login(ctx context.Context) error {
	credentials := url.UserPassword(v.username, v.password)
	err := v.govmomiClient.Login(ctx, credentials)
	if err != nil {
		return err
	}
	return nil
}

func (v *VCenterManager) FindVM(ctx context.Context, inventoryPath string) (*object.VirtualMachine, error) {

	vm, err := v.finder.VirtualMachine(ctx, inventoryPath)
	if err != nil {
		return nil, err
	}

	return vm, nil
}

func (v *VCenterManager) OperationsManager(ctx context.Context, vm *object.VirtualMachine) OpsManager {
	return guest.NewOperationsManager(v.vimClient, vm.Reference())
}

func (v *VCenterManager) GuestManager(ctx context.Context, opsManager OpsManager, username, password string) (*guest_manager.GuestManager, error) {
	processManager, err := opsManager.ProcessManager(ctx)
	if err != nil {
		return nil, err
	}
	auth := types.NamePasswordAuthentication{
		Username: username,
		Password: password,
	}
	return guest_manager.NewGuestManager(auth, processManager), nil
}
