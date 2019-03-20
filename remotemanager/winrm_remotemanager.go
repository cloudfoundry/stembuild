package remotemanager

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
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

func (w *WinRM) CanReachVM() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", w.host, WinrmPort), time.Duration(time.Second*60))
	if err != nil {
		return fmt.Errorf("host %s is unreachable. Please ensure WinRM is enabled and the IP is correct", w.host)
	}
	conn.Close()
	return nil
}

func (w *WinRM) CanLoginVM() error {
	endpoint := winrm.NewEndpoint(w.host, WinrmPort, false, true, nil, nil, nil, time.Second*60)
	winrmClient, err := winrm.NewClient(endpoint, w.username, w.password)
	if err != nil {
		return fmt.Errorf("failed to create winrm client: %s", err)
	}

	s, err := winrmClient.CreateShell()
	if err != nil {
		return errors.New("username and password for given IP is invalid")
	}
	s.Close()

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

	// We override Stderr because WinRM Copy output a lot of XML status messages to StdErr
	// even though they are not errors. In addition, these status messages are difficult to read
	// and add little customer value. WinRM does not have an output override for Copy yet
	reader, tmpStdOut, _ := os.Pipe()
	oldStdErr := os.Stderr
	os.Stderr = tmpStdOut

	defer func() {
		os.Stderr = oldStdErr
		_ = tmpStdOut.Close()
		_ = reader.Close()
	}()

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
