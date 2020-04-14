package remotemanager

import (
	"github.com/masterzen/winrm"
	"time"
)

// todo: test winrm_remotemanager.CanLoginVM() by faking out this factory?
type WinRMClientFactory struct {
	host     string
	username string
	password string
}

func NewWinRmClientFactory(host, username, password string) *WinRMClientFactory {
	return &WinRMClientFactory{host: host, username: username, password: password}
}

func (f *WinRMClientFactory) Build(timeout time.Duration) (WinRMClient, error) {
	// todo run integration tests, then move endpoint creation out
	endpoint := winrm.NewEndpoint(f.host, 5985, false, true, nil, nil, nil, timeout)
	client, err := winrm.NewClient(endpoint, f.username, f.password)
	return client, err
}
