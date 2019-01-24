package iaas_cli_test

import (
	"github.com/cloudfoundry-incubator/stembuild/iaas_cli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GovcCli", func() {
	Describe("RunWithOutputs", func() {
		var runner iaas_cli.GovcRunner

		BeforeEach(func() {
			runner = iaas_cli.GovcRunner{}
		})

		It("lists the devices for a known VCenter VM", func() {
			out, _, err := runner.RunWithOutput([]string{"device.ls", "-vm", "/canada-dc/vm/calgary/stembuild-package-integration-tests-base-vm"})
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(Equal(
				`ide-200            VirtualIDEController          IDE 0
ide-201            VirtualIDEController          IDE 1
ps2-300            VirtualPS2Controller          PS2 controller 0
pci-100            VirtualPCIController          PCI controller 0
sio-400            VirtualSIOController          SIO controller 0
keyboard-600       VirtualKeyboard               Keyboard
pointing-700       VirtualPointingDevice         Pointing device; Device
video-500          VirtualMachineVideoCard       Video card
vmci-12000         VirtualMachineVMCIDevice      Device on the virtual machine PCI bus that provides support for the virtual machine communication interface
lsilogic-sas-1000  VirtualLsiLogicSASController  LSI Logic SAS
ahci-15000         VirtualAHCIController         AHCI
disk-1000-0        VirtualDisk                   41,943,040 KB
floppy-8000        VirtualFloppy                 Remote
ethernet-0         VirtualE1000e                 DVSwitch: a7 fa 3a 50 a9 72 57 5a-56 d1 f3 82 a6 1e 2a ed
cdrom-16000        VirtualCdrom                  Remote device
`,
			))
		})

		It("returns exit code 1, if VM doesn't exist", func() {
			out, exitCode, err := runner.RunWithOutput([]string{"device.ls", "-vm", "/vm/does/not/exist"})
			Expect(err).NotTo(HaveOccurred())
			Expect(exitCode).To(Equal(1))
			Expect(out).To(Equal(""))
		})
	})
})
