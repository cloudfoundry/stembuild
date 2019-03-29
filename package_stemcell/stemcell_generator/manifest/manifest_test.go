package manifest_test

import (
	"bytes"
	"fmt"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/stemcell_generator/manifest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manifest", func() {

	const format = `---
name: bosh-vsphere-esxi-windows%[1]s-go_agent
version: '%[2]s'
sha1: '%[3]s'
operating_system: windows%[1]s
cloud_properties:
  infrastructure: vsphere
  hypervisor: esxi
stemcell_formats:
- vsphere-ovf
- vsphere-ova
`
	Describe("Manifest", func() {
		It("should return a manifest", func() {
			manifestGenerator := manifest.NewManifestGenerator("1709", "1709.999")
			fakeImage := bytes.NewReader([]byte("An image"))

			manifest, err := manifestGenerator.Manifest(fakeImage)

			Expect(err).ToNot(HaveOccurred())

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(manifest)
			Expect(err).NotTo(HaveOccurred())
			s := buf.String()

			//output of `echo "An image" | shasum`
			// TODO: work out where this shasum comes from
			shaSum := "bf8a473a2baa3988b4e7fc4702c35303cdf6df6b"

			Expect(s).To(Equal(fmt.Sprintf(format, "1709", "1709.999", shaSum)))
		})
	})
})
