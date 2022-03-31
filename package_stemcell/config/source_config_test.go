package config_test

import (
	. "github.com/cloudfoundry/stembuild/package_stemcell/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SourceConfig", func() {
	Describe("GetSource", func() {
		It("returns no error when VMDK configured correctly", func() {
			config := SourceConfig{
				Vmdk: "/some/path/to/a/file",
			}
			source, err := config.GetSource()
			Expect(err).NotTo(HaveOccurred())
			Expect(source).To(Equal(VMDK))
		})

		It("returns an error when no configuration provided", func() {
			config := SourceConfig{}
			source, err := config.GetSource()
			Expect(err).To(MatchError("no configuration was provided"))
			Expect(source).To(Equal(NIL))
		})

		It("return no error when vCenter configured correctly", func() {
			config := SourceConfig{
				VmInventoryPath: "/my-datacenter/vm/my-folder/my-vm",
				Username:        "user",
				Password:        "pass",
				URL:             "https://vcenter.test",
			}
			source, err := config.GetSource()
			Expect(err).NotTo(HaveOccurred())
			Expect(source).To(Equal(VCENTER))

		})

		It("returns an error when both VMDK and Vcenter configured", func() {
			config := SourceConfig{
				Vmdk:            "/some/path/to/a/file",
				VmInventoryPath: "/my-datacenter/vm/my-folder/my-vm",
				Username:        "user",
				Password:        "pass",
				URL:             "https://vcenter.test",
			}
			source, err := config.GetSource()
			Expect(err).To(MatchError("configuration provided for VMDK & vCenter sources"))
			Expect(source).To(Equal(NIL))

		})

		It("returns an error when Vcenter configurations only partially specified", func() {
			config := SourceConfig{
				VmInventoryPath: "/my-datacenter/vm/my-folder/my-vm",
				Username:        "user",
				URL:             "https://vcenter.test",
			}
			source, err := config.GetSource()
			Expect(err).To(MatchError("missing vCenter configurations"))
			Expect(source).To(Equal(NIL))
		})

	})

})
