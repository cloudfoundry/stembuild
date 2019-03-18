package manifest

import (
	"bytes"
	"io"
)

type ManifestGenerator struct{}

func (m *ManifestGenerator) Manifest(reader io.Reader) (io.Reader, error) {
	return bytes.NewReader([]byte("Im not nil")), nil
}
