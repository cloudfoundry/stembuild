package remotemanager

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cloudfoundry-incubator/winrmcp/winrmcp"
	"github.com/masterzen/winrm"
)

const WinrmPort = 5985

type WinRM struct {
	host     string
	username string
	password string
}

func NewWinRM(host, username, password string) RemoteManager {
	return &WinRM{host, username, password}
}

func (w *WinRM) CanConnectToVM() error {
	endpoint := winrm.NewEndpoint(w.host, WinrmPort, false, true, nil, nil, nil, time.Second*60)
	winrmClient, err := winrm.NewClient(endpoint, w.username, w.password)
	if err != nil {
		fmt.Printf("Failed to create WinRM client")
		return err
	}

	s, err := winrmClient.CreateShell()
	if err != nil {
		fmt.Printf("Failed to connect to VM")
		return err
	}
	defer s.Close()

	return nil
}

func (w *WinRM) UploadArtifact(sourceFilePath, destinationFilePath string) error {
	client, err := winrmcp.New(w.host, &winrmcp.Config{
		Auth:                  winrmcp.Auth{User: w.username, Password: w.password},
		Https:                 false,
		Insecure:              true,
		OperationTimeout:      time.Second * 60,
		MaxOperationsPerShell: 15,
	})

	if err != nil {
		return err
	}

	return client.Copy(sourceFilePath, destinationFilePath)
}

func (w *WinRM) ExtractArchive(source, destination string) error {
	command := fmt.Sprintf("powershell.exe Expand-Archive %s %s -Force", source, destination)
	err := w.ExecuteCommand(command)
	return err
}

func (w *WinRM) ExecuteCommand(command string) error {
	endpoint := winrm.NewEndpoint(w.host, 5985, false, true, nil, nil, nil, time.Second*60)
	client, err := winrm.NewClient(endpoint, w.username, w.password)
	if err != nil {
		return err
	}
	errBuffer := new(bytes.Buffer)
	exitCode, err := client.Run(command, os.Stdout, io.MultiWriter(errBuffer, os.Stderr))
	if err == nil && exitCode != 0 {
		err = fmt.Errorf("powershell encountered an issue: %s", errBuffer.String())
	}
	return err
}
