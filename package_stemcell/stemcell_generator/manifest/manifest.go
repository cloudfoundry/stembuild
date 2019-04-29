package manifest

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
)

type ManifestGenerator struct {
	os      string
	version string
}

func NewManifestGenerator(os, version string) *ManifestGenerator {
	return &ManifestGenerator{os, version}
}

func (m *ManifestGenerator) Manifest(image io.Reader) (io.Reader, error) {
	const manifestTemplate = `---
name: bosh-vsphere-esxi-windows%[1]s-go_agent
version: '%[2]s'
sha1: '%[3]x'
operating_system: windows%[1]s
cloud_properties:
  infrastructure: vsphere
  hypervisor: esxi
stemcell_formats:
- vsphere-ovf
- vsphere-ova
`
	sha1Hash := sha1.New()

	_, err := io.Copy(sha1Hash, image)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate image shasum: %s", err)
	}

	sum := sha1Hash.Sum(nil)

	manifestContent := fmt.Sprintf(manifestTemplate, m.os, m.version, sum)
	return bytes.NewReader([]byte(manifestContent)), nil
}
