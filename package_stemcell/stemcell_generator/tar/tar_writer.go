package tar

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

type TarWriter struct {
}

func NewTarWriter() *TarWriter {
	return &TarWriter{}
}

//counterfeiter:generate . Tarable
type Tarable interface {
	io.Reader
	Size() int64
	Name() string
}

func (t *TarWriter) Write(filename string, tarables ...Tarable) error {
	// create tar file
	tarFile, _ := os.Create(filename)
	defer tarFile.Close()

	gzw := gzip.NewWriter(tarFile)
	tarfileWriter := tar.NewWriter(gzw)

	for _, t := range tarables {
		// prepare the tar header
		header := new(tar.Header)
		header.Name = t.Name()
		header.Size = t.Size()
		header.Mode = int64(os.FileMode(0644))

		tarfileWriter.WriteHeader(header)
		io.Copy(tarfileWriter, t)
	}
	tarfileWriter.Close()
	gzw.Close()

	return nil
}
