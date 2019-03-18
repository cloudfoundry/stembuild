package manifest_test

import (
	"bytes"
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/stemcell_generator/manifest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manifest", func() {
	Describe("Manifest", func() {
		It("should return a manifest", func() {
			manifestGenerator := &manifest.ManifestGenerator{}
			fakeImage := bytes.NewReader([]byte("An image"))

			manifest, err := manifestGenerator.Manifest(fakeImage)

			Expect(err).ToNot(HaveOccurred())
			Expect(manifest).NotTo(BeNil())
		})
	})
})
