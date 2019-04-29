package filename

import "fmt"

type filenameGenerator struct {
	os      string
	version string
}

func NewFilenameGenerator(os, version string) *filenameGenerator {
	return &filenameGenerator{os: os, version: version}
}

func (f *filenameGenerator) Filename() string {
	return fmt.Sprintf("bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz", f.version, f.os)
}
