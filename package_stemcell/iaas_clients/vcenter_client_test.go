package iaas_clients_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry-incubator/stembuild/iaas_cli/iaas_clifakes"
	. "github.com/cloudfoundry-incubator/stembuild/package_stemcell/iaas_clients"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VcenterClient", func() {
	var (
		runner                  *iaas_clifakes.FakeCliRunner
		username, password, url string
		vcenterClient           *VcenterClient
		credentialUrl           string
	)

	BeforeEach(func() {
		runner = &iaas_clifakes.FakeCliRunner{}
		username, password, url = "username", "password", "url"
		vcenterClient = NewVcenterClient(username, password, url, runner)
		credentialUrl = fmt.Sprintf("%s:%s@%s", username, password, url)
	})

	Context("ValidateCredentials", func() {
		It("When the login credentials are correct, login is successful", func() {
			expectedArgs := []string{"about", "-u", credentialUrl}

			runner.RunReturns(0)
			err := vcenterClient.ValidateCredentials()
			argsForRun := runner.RunArgsForCall(0)

			Expect(err).To(Not(HaveOccurred()))
			Expect(runner.RunCallCount()).To(Equal(1))
			Expect(argsForRun).To(Equal(expectedArgs))
		})

		It("When the login credentials are incorrect, login is a failure", func() {
			expectedArgs := []string{"about", "-u", credentialUrl}

			runner.RunReturns(1)
			err := vcenterClient.ValidateCredentials()
			argsForRun := runner.RunArgsForCall(0)

			Expect(err).To(HaveOccurred())
			Expect(runner.RunCallCount()).To(Equal(1))
			Expect(argsForRun).To(Equal(expectedArgs))
			Expect(err).To(MatchError("vcenter_client - invalid credentials for: url"))
		})
	})

	Context("validateUrl", func() {
		It("When the url is valid, there is no error", func() {
			expectedArgs := []string{"about", "-u", url}

			runner.RunReturns(0)
			err := vcenterClient.ValidateUrl()
			argsForRun := runner.RunArgsForCall(0)

			Expect(err).To(Not(HaveOccurred()))
			Expect(runner.RunCallCount()).To(Equal(1))
			Expect(argsForRun).To(Equal(expectedArgs))
		})

		It("When the url is invalid, there is an error", func() {
			expectedArgs := []string{"about", "-u", url}

			runner.RunReturns(1)
			err := vcenterClient.ValidateUrl()
			argsForRun := runner.RunArgsForCall(0)

			Expect(err).To(HaveOccurred())
			Expect(runner.RunCallCount()).To(Equal(1))
			Expect(argsForRun).To(Equal(expectedArgs))
			Expect(err).To(MatchError("vcenter_client - unable to validate url: url"))
		})
	})

	Context("FindVM", func() {
		It("If the VM path is valid, and the VM is found", func() {
			expectedArgs := []string{"find", "-maxdepth=0", "-u", credentialUrl, "validVMPath"}
			runner.RunReturns(0)
			err := vcenterClient.FindVM("validVMPath")
			argsForRun := runner.RunArgsForCall(0)

			Expect(err).To(Not(HaveOccurred()))
			Expect(runner.RunCallCount()).To(Equal(1))
			Expect(argsForRun).To(Equal(expectedArgs))
		})

		It("If the VM path is invalid", func() {
			expectedArgs := []string{"find", "-maxdepth=0", "-u", credentialUrl, "invalidVMPath"}
			runner.RunReturns(1)
			err := vcenterClient.FindVM("invalidVMPath")
			argsForRun := runner.RunArgsForCall(0)

			Expect(err).To(HaveOccurred())
			Expect(runner.RunCallCount()).To(Equal(1))
			Expect(argsForRun).To(Equal(expectedArgs))
			Expect(err).To(MatchError("vcenter_client - unable to find VM: invalidVMPath. Ensure your inventory path is formatted properly and includes \"vm\" in its path, example: /my-datacenter/vm/my-folder/my-vm-name"))
		})
	})

	Describe("RemoveDevice", func() {
		It("Removes a device from the given VM", func() {
			runner.RunReturns(0)
			err := vcenterClient.RemoveDevice("validVMPath", "device")

			Expect(err).To(Not(HaveOccurred()))
			expectedArgs := []string{"device.remove", "-u", credentialUrl, "-vm", "validVMPath", "device"}
			Expect(runner.RunArgsForCall(0)).To(Equal(expectedArgs))
		})

		It("Returns an error if VCenter reports a failure removing a device", func() {
			runner.RunReturns(1)
			err := vcenterClient.RemoveDevice("VMPath", "deviceName")

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("vcenter_client - deviceName could not be removed"))
		})
	})

	Describe("EjectCDrom", func() {
		It("Ejects a cd rom from the given VM", func() {
			runner.RunReturns(0)
			err := vcenterClient.EjectCDRom("validVMPath", "deviceName")

			Expect(err).To(Not(HaveOccurred()))
			expectedArgs := []string{"device.cdrom.eject", "-u", credentialUrl, "-vm", "validVMPath", "-device", "deviceName"}
			Expect(runner.RunArgsForCall(0)).To(Equal(expectedArgs))
		})

		It("Returns an error if VCenter reports a failure ejecting the cd rom", func() {
			runner.RunReturns(1)
			err := vcenterClient.EjectCDRom("VMPath", "deviceName")

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("vcenter_client - deviceName could not be ejected"))
		})
	})

	Context("ListDevices", func() {
		var govcListDevicesOutput = `ide-200            VirtualIDEController          IDE 0
ide-201            VirtualIDEController          IDE 1
ps2-300            VirtualPS2Controller          PS2 controller 0
pci-100            VirtualPCIController          PCI controller 0
sio-400            VirtualSIOController          SIO controller 0
floppy-8000        VirtualFloppy                 Remote
ethernet-0         VirtualE1000e                 DVSwitch: a7 fa 3a 50 a9 72 57 5a-56 d1 f3 82 a6 1e 2a ed
`

		It("returns a list of devices for the given VM", func() {
			runner.RunWithOutputReturns(govcListDevicesOutput, 0, nil)

			devices, err := vcenterClient.ListDevices("/path/to/vm")

			Expect(err).NotTo(HaveOccurred())
			Expect(devices).To(ConsistOf(
				"ide-200", "ide-201", "ps2-300", "pci-100", "sio-400", "floppy-8000", "ethernet-0",
			))

			Expect(runner.RunWithOutputArgsForCall(0)).To(Equal([]string{"device.ls", "-vm", "/path/to/vm"}))
		})

		It("returns an error if govc runner returns non zero exit code", func() {
			runner.RunWithOutputReturns("", 1, nil)

			_, err := vcenterClient.ListDevices("/path/to/vm")

			Expect(err).To(MatchError("vcenter_client - failed to list devices in vCenter, govc exit code 1"))
		})

		It("returns an error if RunWithOutput encounters an error", func() {
			runner.RunWithOutputReturns("", 0, errors.New("some environment error"))

			_, err := vcenterClient.ListDevices("/path/to/vm")

			Expect(err).To(MatchError("vcenter_client - failed to parse list of devices. Err: some environment error"))
		})

		It("returns govc exit code error, when both govc exit code is non zero and RunWithOutput encounters an error", func() {
			runner.RunWithOutputReturns("", 1, errors.New("some environment error"))

			_, err := vcenterClient.ListDevices("/path/to/vm")

			Expect(err).To(MatchError("vcenter_client - failed to list devices in vCenter, govc exit code 1"))
		})
	})

	Context("ExportVM", func() {
		var destinationDir string
		BeforeEach(func() {
			destinationDir, _ = ioutil.TempDir(os.TempDir(), "destinationDir")
		})
		It("exports the VM to local machine from vcenter using vm inventory path", func() {
			expectedArgs := []string{"export.ovf", "-u", credentialUrl, "-sha", "1", "-vm", "validVMPath", destinationDir}
			runner.RunReturns(0)
			err := vcenterClient.ExportVM("validVMPath", destinationDir)

			Expect(err).To(Not(HaveOccurred()))
			Expect(runner.RunCallCount()).To(Equal(1))

			argsForRun := runner.RunArgsForCall(0)
			Expect(argsForRun).To(Equal(expectedArgs))
		})

		It("Returns an error message if ExportVM fails to export the VM", func() {
			vmInventoryPath := "validVMPath"
			expectedArgs := []string{"export.ovf", "-u", credentialUrl, "-sha", "1", "-vm", vmInventoryPath, destinationDir}
			runner.RunReturns(1)
			err := vcenterClient.ExportVM("validVMPath", destinationDir)

			expectedErrorMsg := fmt.Sprintf("vcenter_client - %s could not be exported", vmInventoryPath)
			Expect(err).To(HaveOccurred())
			Expect(runner.RunCallCount()).To(Equal(1))

			argsForRun := runner.RunArgsForCall(0)
			Expect(argsForRun).To(Equal(expectedArgs))
			Expect(err.Error()).To(Equal(expectedErrorMsg))
		})

		It("prints an appropriate error message if the given directory doesn't exist", func() {
			err := vcenterClient.ExportVM("validVMPath", "/FooBar/stuff")
			Expect(err).To(HaveOccurred())

			Expect(err.Error()).To(Equal("vcenter_client - provided destination directory: /FooBar/stuff does not exist"))
		})
	})
})
