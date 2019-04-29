package stemcell_generator

import (
	"fmt"
	"io"
)

//go:generate counterfeiter . ManifestGenerator
type ManifestGenerator interface {
	Manifest(reader io.Reader) (io.Reader, error)
}

//go:generate counterfeiter . FileNameGenerator
type FileNameGenerator interface {
	FileName() string
}

//go:generate counterfeiter . TarWriter
type TarWriter interface {
	Write(string, ...io.Reader) error
}
type StemcellGenerator struct {
	manifestGenerator ManifestGenerator
	fileNameGenerator FileNameGenerator
	tarWriter         TarWriter
}

func NewStemcellGenerator(m ManifestGenerator, f FileNameGenerator, t TarWriter) *StemcellGenerator {
	return &StemcellGenerator{m, f, t}
}

func (g *StemcellGenerator) Generate(image io.Reader) error {
	manifest, err := g.manifestGenerator.Manifest(image)
	if err != nil {
		return fmt.Errorf("failed to generate stemcell manifest: %s", err)
	}
	filename := g.fileNameGenerator.FileName()

	err = g.tarWriter.Write(filename, image, manifest)
	if err != nil {
		return fmt.Errorf("failed to generate stemcell tarball: %s", err)
	}

	return nil
}
