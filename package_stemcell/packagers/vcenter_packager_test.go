package packagers

import (
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VcenterPackager", func() {
	Context("ValidateSourceParameters", func() {
		It("returns an error if the vCenter url is invalid", func() {
			sourceConfig := config.SourceConfig{Password: "password", URL: "http://invalid.url", Username: "username", VmInventoryPath: "/valid-datacenter/vm/valid/path"}
			client := FakeVcenterClient{Username: sourceConfig.Username, Password: sourceConfig.Password, Url: sourceConfig.URL, InvalidUrl: true}
			packager := VCenterPackager{SourceConfig: sourceConfig, Client: client}
			err := packager.ValidateSourceParameters()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("please provide a valid vCenter URL"))

		})
		It("returns an error if the vCenter credentials are not valid", func() {
			sourceConfig := config.SourceConfig{Password: "invalidPassword", URL: "http://foocenter.bar", Username: "username", VmInventoryPath: "/valid-datacenter/vm/valid/path"}
			client := FakeVcenterClient{Username: sourceConfig.Username, Password: sourceConfig.Password, Url: sourceConfig.URL, InvalidCredentials: true}
			packager := VCenterPackager{SourceConfig: sourceConfig, Client: client}

			err := packager.ValidateSourceParameters()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("please provide valid credentials for"))
		})

		It("returns an error if VM given does not exist ", func() {
			sourceConfig := config.SourceConfig{Password: "password", URL: "http://foocenter.bar", Username: "username", VmInventoryPath: "/invalid-datacenter/invalid/path"}
			client := FakeVcenterClient{Username: sourceConfig.Username, Password: sourceConfig.Password, Url: sourceConfig.URL, InvalidVmInventoryPath: true}
			packager := VCenterPackager{SourceConfig: sourceConfig, Client: client}

			err := packager.ValidateSourceParameters()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("VM path is invalid\nPlease make sure to format your inventory path correctly using the 'vm' keyword. Example: /my-datacenter/vm/my-folder/my-vm-name"))
		})
		It("returns no error if all source parameters are valid", func() {
			sourceConfig := config.SourceConfig{Password: "password", URL: "http://foocenter.bar", Username: "username", VmInventoryPath: "/valid-datacenter/vm/valid/path"}
			client := FakeVcenterClient{Username: sourceConfig.Username, Password: sourceConfig.Password, Url: sourceConfig.URL}
			packager := VCenterPackager{SourceConfig: sourceConfig, Client: client}

			err := packager.ValidateSourceParameters()

			Expect(err).NotTo(HaveOccurred())
		})
	})
})
