package remotemanager

import (
	"bytes"
	"fmt"
	"github.com/masterzen/winrm"
	"io"
	"net"
	"os"
	"time"

	"github.com/cloudfoundry-incubator/winrmcp/winrmcp"
)

const WinRmPort = 5985
const WinRmTimeout = 60 * time.Second

type WinRM struct {
	host          string
	username      string
	password      string
	clientFactory WinRMClientFactoryI
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . WinRMClient
type WinRMClient interface {
	Run(command string, stdout io.Writer, stderr io.Writer) (int, error)
	CreateShell() (*winrm.Shell, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . WinRMClientFactoryI
type WinRMClientFactoryI interface {
	Build(timeout time.Duration) (WinRMClient, error)
}

func NewWinRM(host string, username string, password string, clientFactory WinRMClientFactoryI) RemoteManager {
	return &WinRM{host, username, password, clientFactory}
}

func (w *WinRM) CanReachVM() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", w.host, WinRmPort), time.Duration(time.Second*60))
	if err != nil {
		return fmt.Errorf("host %s is unreachable. Please ensure WinRM is enabled and the IP is correct: %s", w.host, err)
	}
	conn.Close()
	return nil
}

func (w *WinRM) CanLoginVM() error {
	winrmClient, err := w.clientFactory.Build(WinRmTimeout)

	if err != nil {
		return fmt.Errorf("failed to create winrm client: %s", err)
	}

	s, err := winrmClient.CreateShell()
	if err != nil {
		return fmt.Errorf("failed to create winrm shell: %s", err)
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
	_, err := w.ExecuteCommand(command)
	return err
}

func (w *WinRM) ExecuteCommandWithTimeout(command string, timeout time.Duration) (int, error) {
	client, err := w.clientFactory.Build(timeout)
	if err != nil {
		return -1, err
	}
	errBuffer := new(bytes.Buffer)
	exitCode, err := client.Run(command, os.Stdout, io.MultiWriter(errBuffer, os.Stderr))
	if err == nil && exitCode != 0 {
		err = fmt.Errorf("%s: %s", PowershellExecutionErrorMessage, errBuffer.String())
	}
	return exitCode, err
}

func (w *WinRM) ExecuteCommand(command string) (int, error) {
	exitCode, err := w.ExecuteCommandWithTimeout(command, WinRmTimeout)
	return exitCode, err
}
