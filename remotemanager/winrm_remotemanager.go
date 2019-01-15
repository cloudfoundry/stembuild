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

type WinRM struct {
	host     string
	username string
	password string
}

func NewWinRM(host, username, password string) RemoteManager {
	return &WinRM{host, username, password}
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
	command := fmt.Sprintf("powershell.exe Expand-Archive %s %s", source, destination)
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
	exitCode, err := client.RunWithInput(command, os.Stdout, io.MultiWriter(errBuffer, os.Stderr), os.Stdin)
	if err == nil && exitCode != 0 {
		err = fmt.Errorf("powershell encountered an issue: %s", errBuffer.String())
	}
	return err
}
