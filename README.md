# Stembuild

The stembuild binary is used to build BOSH stemcells for **Windows 2012R2**,**Windows Server, version v1709**, **Windows Server, version 1803**, **Windows Server 2019** on **vSphere**. 

**Instructions**: See [here](https://bosh.io/docs/windows-stemcell-create/) for instructions to build Windows stemcells for vSphere.

## Installation
Download the latest stembuild from the [Releases Page](https://github.com/cloudfoundry/stembuild/releases) that corresponds to the operating system of your local host and the stemcell version that you want to build

## Dependencies
[LGPO](https://www.microsoft.com/en-us/download/details.aspx?id=55319) must be downloaded in the same folder as your `stembuild`

## Current Commands
```
stembuild version <STEMCELL-VERSION>, Windows Stemcell Building Tool

Usage: stembuild <global options> <command> <command flags>

Commands:
  help		Describe commands and their syntax
  package	Create a BOSH Stemcell from a VMDK file or a provisioned vCenter VM
  construct	Provisions and syspreps an existing VM on vCenter, ready to be packaged into a stemcell

Global Options:
  -color	Colorize debug output
  -debug	Print lots of debugging information
  -v		Stembuild version (shorthand)
  -version	Show Stembuild version

```
## `stembuild construct`

This command provisions and syspreps an existing VM on vCenter. It prepares a VM to be used by `stembuild package`.

```
stembuild construct -vm-ip <IP of VM> -vm-username <vm username> -vm-password <vm password>  -vcenter-url <vCenter URL> -vcenter-username <vCenter username> -vcenter-password <vCenter password> -vm-inventory-path <vCenter VM inventory path>
```

### Requirements
- LGPO.zip in current working directory
- Running Windows VM with:
	- Up-to-date Operating System
	- Reachable by IP over port 5985
	- Username and password with Administrator privileges
	- vCenter URL, username and password
	- vCenter Inventory Path
- The `vm-ip`, `vm-username`, `vm-password`, `vcenter-url`, `vcenter-username`, `vcenter-password`, `vm-inventory-path` must be specified

```
Example:
	stembuild construct -vm-ip '10.0.0.5' -vm-username Admin -vm-password 'password' -vcenter-url vcenter.example.com -vcenter-username root -vcenter-password 'password' -vm-inventory-path '/datacenter/vm/folder/vm-name'

Flags:
  -vcenter-ca-certs string
    	filepath for custom ca certs
  -vcenter-password string
    	vCenter password
  -vcenter-url string
    	vCenter url
  -vcenter-username string
    	vCenter username
  -vm-inventory-path string
    	vCenter VM inventory path. (e.g: /<datacenter>/vm/<vm-folder>/<vm-name>)
  -vm-ip string
    	IP of target machine
  -vm-password string
    	Password of target machine. Needs to be wrapped in single quotations.
  -vm-username string
    	Username of target machine
	
```

### Troubleshooting
After running `stembuild construct`, you may find yourself with a connection issue to the VM
- Confirm port 5985 is reachable via something like `nmap [vm-ip] -Pn`


## `stembuild package`

This command creates a BOSH Stemcell from a provisioned vCenter VM 

```
  stembuild package -vcenter-url <vCenter URL> -vcenter-username <vCenter username> -vcenter-password <vCenter password> -vm-inventory-path <vCenter VM inventory path> [-patch-version <patch version string>]
```

*Requirements*:
- VM provisioned using the stembuild construct command
- Access to vCenter environment
- The `vcenter-url`, `vcenter-username`, `vcenter-password`, and `vm-inventory-path` flags must be specified.
- **NOTE**: The 'vm' keyword must be included between the datacenter name and folder name for the vm-inventory-path (e.g: /<datacenter>/vm/<vm-folder>/<vm-name>)
 
```
Example:
 stembuild package -vcenter-url vcenter.example.com -vcenter-username root -vcenter-password 'password' -vm-inventory-path '/my-datacenter/vm/my-folder/my-vm'

Flags:
  -o string
    	Output directory (shorthand)
  -outputDir string
    	Output directory, default is the current working directory.
  -vcenter-ca-certs string
    	filepath for custom ca certs
  -vcenter-password string
    	vCenter password
  -vcenter-url string
    	vCenter url
  -vcenter-username string
    	vCenter username
  -vm-inventory-path string
    	vCenter VM inventory path. (e.g: /<datacenter>/vm/<vm-folder>/<vm-name>)
  -patch-version string
  	Number or name of the patch version for the stemcell being built (e.g: for 2019.12.3 the string would be “3”)

```

### Compiling & Running Stembuild Locally

Assuming you've followed [these instructions](https://bosh.io/docs/windows-stemcell-create/) and you've created a Windows VM at 10.9.9.115 whose Administrator's password is "c1oudc0w".

```bash
    export TARGET_VM_PASSWORD=c1oudc0w VCENTER_PASSWORD='Admin!23'
BOSH_AGENT_REPO=~/workspace/bosh-agent STEMBUILD_VERSION=2019.2 make
GOVC_INSECURE=true out/stembuild -debug \
  construct \
  -vm-ip 10.9.9.115 \
  -vm-username Administrator \
  -vm-password $TARGET_VM_PASSWORD \
  -vcenter-url vcenter-70.nono.io \
  -vcenter-username administrator@vsphere.local \
  -vcenter-password $VCENTER_PASSWORD \
  -vm-inventory-path "/dc/vm/Discovered virtual machine/w2019-stemcell"
GOVC_INSECURE=true out/stembuild -debug \
  package \
  -vcenter-url vcenter-70.nono.io \
  -vcenter-username administrator@vsphere.local \
  -vcenter-password $VCENTER_PASSWORD \
  -vm-inventory-path "/dc/vm/Discovered virtual machine/w2019-stemcell"
```

## [DEPRECATED] Package a Windows Stemcell from a VMDK using `stembuild package`

This command converts a VMDK into a bosh-deployable Windows Stemcell 

The VMware 'ovftool' binary must be on your path or Fusion/Workstation must be installed (both include the 'ovftool').

```
stembuild package -vmdk <path-to-vmdk>
```

*Requirements*
- The VMware 'ovftool' binary must be on your path or Fusion/Workstation must be installed (both include the 'ovftool').
- The `vmdk` flag must be specified.  If the `output` flag is not specified the stemcell will be created in the current working directory.

```
Example:
	stembuild package -vmdk my-1803-vmdk.vmdk
	
	Will create an Windows 1803 stemcell using [vmdk] 'my-1803-vmdk.vmdk'
	The final stemcell will be found in the current working directory.

Flags:
  -o string
    	Output directory (shorthand)
  -outputDir string
    	Output directory, default is the current working directory.
  -vmdk string
    	VMDK file to create stemcell from

```

Process can take between 10 and 20 minutes. See Progress with `-debug` flag.

## Testing

### Testing stembuild itself

Older tests were written using the default testing framework.  However, more recent code
has been test-driven with Ginkgo.  We recommend that any new code be test-driven using Ginkgo.
Below are steps to run the tests:

Make puts some files in `out` dir. To clean state of this dir:
```
make clean
```
To run only unit tests:
```
make units
```
To run integration tests:
```
make integration
```

### Testing stemcell-automation

`stemcell-automation` contains powershell scripts which stembuild runs on the target VM. This directory also contains tests
for the scripts, which are written with the test framework [Pester](https://pester.dev/). These require running in powershell on a Windows environment.

```
cd stemcell-automation
invoke-pester
```

## Vendoring

Vendoring for this project is done using `dep`. 
To sync all the dependencies run
```
dep ensure
```

To add a new dependency run 
```
dep ensure -add <git package url>
```
like 
```
dep ensure -add github.com/google/subcommands
```

To check if dependencies are in sync or not run
```
dep sync
```
The output should be nothing if there are no out-of-sync dependencies.


## Compile stembuild for macOS

You can use stembuild on macOS by following the below steps:

- Download or clone the bosh-agent repository
```
git clone https://github.com/cloudfoundry/bosh-agent.git
```
- Download or clone the stembuild repository and navigate to it
```
git clone https://github.com/cloudfoundry/stembuild.git
cd stembuild
```
- Download the latest released artifact from [Stemcell Automation GitHub Repo](https://github.com/cloudfoundry-incubator/bosh-windows-stemcell-automation/releases)
- Use `make build` to build stembuild for macOS, providing the corresponding values for the bosh-agent path, Stemcell Version you would like to build, and the Stemcell Automation package you downloaded
```
BOSH_AGENT_REPO=PathToBoshAgentRepo make STEMCELL_VERSION=StemcellVersion AUTOMATION_PATH=PathToStemcellAutomationZip build
```
