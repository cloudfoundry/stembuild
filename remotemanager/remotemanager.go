package remotemanager

import "time"

//go:generate counterfeiter . RemoteManager

const PowershellExecutionErrorMessage = "powershell encountered an issue"

type RemoteManager interface {
	UploadArtifact(source, destination string) error
	ExtractArchive(source, destination string) error
	ExecuteCommand(command string) (int, error)
	ExecuteCommandWithTimeout(command string, timeout time.Duration) (int, error)
	CanReachVM() error
	CanLoginVM() error
}
