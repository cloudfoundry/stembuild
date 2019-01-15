package remotemanager

//go:generate mockgen -source=remotemanager.go -destination=mock/mock_remotemanager.go RemoteManager

type RemoteManager interface {
	UploadArtifact(source, destination string) error
	ExtractArchive(source, destination string) error
	ExecuteCommand(file string) error
}
