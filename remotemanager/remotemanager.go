package remotemanager

//go:generate counterfeiter . RemoteManager

type RemoteManager interface {
	UploadArtifact(source, destination string) error
	ExtractArchive(source, destination string) error
	ExecuteCommand(command string) (int, error)
	CanReachVM() error
	CanLoginVM() error
}
