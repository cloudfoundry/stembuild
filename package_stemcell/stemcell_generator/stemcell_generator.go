package stemcell_generator

import "io"

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
	Write(string, ...io.Reader)
}
type StemcellGenerator struct {
	manifestGenerator ManifestGenerator
	fileNameGenerator FileNameGenerator
	tarWriter TarWriter
}

func NewStemcellGenerator(m ManifestGenerator, f FileNameGenerator, t TarWriter) *StemcellGenerator {
	return &StemcellGenerator{m, f, t}
}

func (g *StemcellGenerator) Generate(image io.Reader) error {
	manifest, _ := g.manifestGenerator.Manifest(image)
	filename := g.fileNameGenerator.FileName()
	g.tarWriter.Write(filename, image, manifest)
	return nil
}
