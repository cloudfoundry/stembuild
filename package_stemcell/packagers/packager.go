package packagers

import "io"

type Packager struct {
}

//go:generate counterfeiter . Source
type Source interface {
	ArtifactReader() io.Reader
}
func (p *Packager) Package() error {
	return nil
}

