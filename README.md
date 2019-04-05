# Stembuild

The stembuild binary is used to build BOSH stemcells for **Windows 2012R2**,**Windows Server, version v1709** and **Windows Server, version 1803** on **vSphere**. See [here](https://github.com/cloudfoundry-incubator/bosh-windows-stemcell-builder/wiki/Creating-a-vSphere-Stemcell-by-Hand) for instructions to build Windows stemcells for vSphere.

It can convert a prepared .vmdk into a stemcell with the appropriate metadata.

Download the latest stembuild from the [Releases page](https://github.com/cloudfoundry-incubator/stembuild/releases).

## Dependencies
The VMware 'ovftool' binary must be on your path or Fusion/Workstation must be installed (both include the 'ovftool').

## Installation

To install `stembuild` go to [Releases](https://github.com/cloudfoundry-incubator/stembuild/releases)

## Current Commands
```
stembuild version 0.21.45, Windows Stemcell Building Tool

Usage: stembuild <global options> <command> <command args>

Commands:
  help		Describe commands and their syntax
  package	Create a BOSH Stemcell from a VMDK file

Global Options:
  -color	Colorize debug output
  -debug	Print lots of debugging information
  -v		Stembuild version (shorthand)
  -version	Show Stembuild version

```
## Create a Windows Stemcell from a VMDK

This command converts a VMDK into a bosh-deployable Windows Stemcell 
```
stembuild package -vmdk <path-to-vmdk>

Create a BOSH Stemcell from a VMDK file

The [vmdk], [stemcellVersion], and [os] flags must be specified.  If the [output] flag is
not specified the stemcell will be created in the current working directory.

Requirements:
	The VMware 'ovftool' binary must be on your path or Fusion/Workstation
	must be installed (both include the 'ovftool').

Examples:
	stembuild package -vmdk disk.vmdk

	Will create an Windows 1803 stemcell using [vmdk] 'disk.vmdk', and set the stemcell version to 1.2.
	The final stemcell will be found in the current working directory.

Flags:
  -o string
    	Output directory (shorthand)
  -outputDir string
    	Output directory, default is the current working directory.
  -s string
    	Stemcell version (shorthand)
  -vmdk string
    	VMDK file to create stemcell from
      
```

Process can take between 10 and 20 minutes. See Progress with `-debug` flag.

## Testing

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


## Vendoring

Vendoring for this project is done using `dep`. 
To sync all the dependecies run
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
The output should be nothing if there no out-of-sync dependencies.
