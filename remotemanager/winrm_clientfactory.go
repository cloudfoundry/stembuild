package remotemanager

import (
	"github.com/masterzen/winrm"
	"time"
)

type WinRMClientFactory struct {
	host     string
	port     int
	username string
	password string
}

func NewWinRmClientFactory(host string, port int, username string, password string) *WinRMClientFactory {
	return &WinRMClientFactory{host: host, port: port, username: username, password: password}
}

func (f *WinRMClientFactory) Build(timeout time.Duration) (WinRMClient, error) {
	endpoint := winrm.NewEndpoint(f.host, f.port, false, true, nil, nil, nil, timeout)
	client, err := winrm.NewClient(endpoint, f.username, f.password)
	return client, err
}
