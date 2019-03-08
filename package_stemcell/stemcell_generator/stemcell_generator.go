package stemcell_generator

import "io"

//go:generate counterfeiter . ManifestGenerator
type ManifestGenerator interface {
	Generate(reader io.Reader) (io.Reader, error)
}

type StemcellGenerator struct {
	manifestGenerator ManifestGenerator

}

func NewStemcellGenerator(m ManifestGenerator) *StemcellGenerator {
	return &StemcellGenerator{m}
}

func (g *StemcellGenerator) Generate(image io.Reader) error {
	g.manifestGenerator.Generate(image)
	return nil
}
