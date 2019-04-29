package manifest_test

import (
	"bytes"
	"errors"
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

			//output of `echo -n "An image" | shasum`
			shaSum := "bf8a473a2baa3988b4e7fc4702c35303cdf6df6b"

			Expect(s).To(Equal(fmt.Sprintf(format, "1709", "1709.999", shaSum)))
		})

		It("should return an error if the image returns an error during read", func() {
			manifestGenerator := manifest.NewManifestGenerator("1709", "1709.999")
			fakeImage := FailReader{}
			_, err := manifestGenerator.Manifest(fakeImage)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to calculate image shasum: failed read"))
		})
	})
})

type FailReader struct {
}

func (f FailReader) Read(p []byte) (int, error) {
	return 0, errors.New("failed read")
}
