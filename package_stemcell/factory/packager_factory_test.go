package factory_test

import (
	"github.com/cloudfoundry/stembuild/package_stemcell/config"
	"github.com/cloudfoundry/stembuild/package_stemcell/factory"
	"github.com/cloudfoundry/stembuild/package_stemcell/packagers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Factory", func() {

	outputConfig := config.OutputConfig{
		Os:              "2012R2",
		StemcellVersion: "1200.00",
		OutputDir:       "/tmp/outputDir",
	}

	var packagerFactory *factory.PackagerFactory

	BeforeEach(func() {
		packagerFactory = &factory.PackagerFactory{}
	})

	Describe("GetPackager", func() {
		Context("When VMDK is specified and no vCenter credentials are given", func() {
			It("returns a VMDK packager with no error", func() {
				sourceConfig := config.SourceConfig{
					Vmdk: "path/to/a/vmdk",
				}

				actualPackager, err := packagerFactory.Packager(sourceConfig, outputConfig, 0, false)
				Expect(err).NotTo(HaveOccurred())

				Expect(actualPackager).To(BeAssignableToTypeOf(packagers.VmdkPackager{}))
				Expect(actualPackager).NotTo(BeAssignableToTypeOf(packagers.VCenterPackager{}))
			})
		})

		Context("When all vCenter credentials are given and no VMDK is specified", func() {
			It("returns a vCenter packager with no error", func() {
				sourceConfig := config.SourceConfig{
					Username:        "user",
					Password:        "pass",
					URL:             "some-url",
					VmInventoryPath: "some-vm-inventory-path",
				}

				actualPackager, err := packagerFactory.Packager(sourceConfig, outputConfig, 0, false)
				Expect(err).NotTo(HaveOccurred())

				Expect(actualPackager).To(BeAssignableToTypeOf(packagers.VCenterPackager{}))
				Expect(actualPackager).NotTo(BeAssignableToTypeOf(packagers.VmdkPackager{}))
			})
		})

		Context("When at least one vCenter configuration and VMDK are both specified", func() {
			It("returns an error", func() {
				sourceConfig := config.SourceConfig{
					Vmdk:            "path/to/a/vmdk",
					VmInventoryPath: "some-vm",
				}

				packager, err := packagerFactory.Packager(sourceConfig, outputConfig, 0, false)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("configuration provided for VMDK & vCenter sources"))
				Expect(packager).To(BeNil())
			})
		})

		Context("When partial vCenter credentials are given and no VMDK is specified", func() {
			It("returns an error", func() {
				sourceConfig := config.SourceConfig{
					VmInventoryPath: "some-vm",
					Password:        "pass",
					URL:             "some-url",
				}

				packager, err := packagerFactory.Packager(sourceConfig, outputConfig, 0, false)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("missing vCenter configurations"))
				Expect(packager).To(BeNil())
			})
		})

		Context("When no configuration has been provided", func() {
			It("returns an error", func() {
				sourceConfig := config.SourceConfig{}

				packager, err := packagerFactory.Packager(sourceConfig, outputConfig, 0, false)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("no configuration was provided"))
				Expect(packager).To(BeNil())
			})
		})
	})
})
