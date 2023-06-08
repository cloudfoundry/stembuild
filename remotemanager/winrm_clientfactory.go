package remotemanager

import (
	"time"

	"github.com/masterzen/winrm"
)

type WinRMClientFactory struct {
	host     string
	username string
	password string
}

func NewWinRmClientFactory(host, username, password string) *WinRMClientFactory {
	return &WinRMClientFactory{host: host, username: username, password: password}
}

func (f *WinRMClientFactory) Build(timeout time.Duration) (WinRMClient, error) {
	endpoint := winrm.NewEndpoint(f.host, WinRmPort, false, true, nil, nil, nil, timeout)
	client, err := winrm.NewClient(endpoint, f.username, f.password)
	return client, err
}
