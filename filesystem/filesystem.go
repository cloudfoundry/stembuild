package filesystem

//go:generate mockgen -source=filesystem.go -destination=mock/mock_filesystem.go FileSystem

type FileSystem interface {
	GetAvailableDiskSpace(path string) (uint64, error)
}
