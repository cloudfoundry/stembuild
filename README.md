# Stembuild

The stembuild binary is used to build BOSH stemcells for Windows 2012R2 and 2016 v1709 on vSphere. See [here](https://github.com/cloudfoundry-incubator/bosh-windows-stemcell-builder/wiki/Creating-a-vSphere-Stemcell-by-Hand) for instructions to build Windows stemcells for vSphere.

It can convert a prepared .vmdk into a stemcell with the appropriate metadata.

Download the latest stembuild from the [Releases page](https://github.com/pivotal-cf-experimental/stembuild/releases).

## Dependencies
The VMware 'ovftool' binary must be on your path or Fusion/Workstation must be installed (both include the 'ovftool').

## Create a stemcell from a vmdk

Usage `stembuild-linux [OPTIONS...] [-vmdk FILENAME]
                                  [-output DIRNAME] [-version STEMCELL_VERSION]
                                  [-os OS_VERSION]`

Process takes between 10 and 20 minutes. See Progress with `-debug` flag.

## Create a stemcell from a vhd

Usage `stembuild-linux [-output DIRNAME] apply-patch <patch manifest yml>`

will generate a stemcell in the current working directory

Other options available with `stembuild -h`.

## Compilation

See the [wiki](https://github.com/pivotal-cf-experimental/stembuild/wiki/build-stembuild)

## Testing

Older tests were written using the default testing framework.  However, more recent code
has been test-driven with Ginkgo.  We recommend that any new code be test-driven using Ginkgo.

## Vendoring

Vendoring for this project is done using `dep`.
