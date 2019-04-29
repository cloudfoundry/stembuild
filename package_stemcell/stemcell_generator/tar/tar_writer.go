package tar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

type TarWriter struct {
}

func NewTarWriter() *TarWriter {
	return &TarWriter{}
}

func (t *TarWriter) Write(filename string, readers ...io.Reader) error {

	// create tar file
	tarFile, _ := os.Create(filename)
	defer tarFile.Close()

	gzw := gzip.NewWriter(tarFile)
	tarfileWriter := tar.NewWriter(gzw)

	for index, r := range readers {

		// prepare the tar header
		header := new(tar.Header)
		header.Name = fmt.Sprintf("File%d", index)
		header.Size = 13
		header.Mode = int64(os.ModePerm)

		tarfileWriter.WriteHeader(header)
		io.Copy(tarfileWriter, r)
	}
	tarfileWriter.Close()
	gzw.Close()

	return nil
}
