package packagers

import (
	"fmt"
	"io"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

type Packager struct {
	source            Source
	stemcellGenerator StemcellGenerator
}

//counterfeiter:generate . Source
type Source interface {
	ArtifactReader() (io.Reader, error)
}

//counterfeiter:generate . StemcellGenerator
type StemcellGenerator interface {
	Generate(reader io.Reader) error
}

func NewPackager(s Source, g StemcellGenerator) *Packager {
	return &Packager{source: s, stemcellGenerator: g}
}

func (p *Packager) Package() error {
	artifact, err := p.source.ArtifactReader()
	if err != nil {
		return fmt.Errorf("packager failed to retrieve artifact: %s", err)
	}
	err = p.stemcellGenerator.Generate(artifact)
	if err != nil {
		return fmt.Errorf("packager failed to generate stemcell: %s", err)
	}

	return nil
}
