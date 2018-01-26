# Instructions

If you have a `.vhd` file and a `.patch` file for that VHD, then:

```
stembuild -vhd my-vhd.vhd -delta patchfile.patch -version 1000.0
```

will create a stemcell with version `1000.0` in your current working directory.

Process takes between 10 and 20 minutes. See Progress with `-debug` flag.

Other options available with `stembuild -h`.

# Compilation

See the [wiki](https://github.com/pivotal-cf-experimental/stembuild/wiki/build-stembuild)

# Testing

Older tests were written using the default testing framework.  However, more recent code
has been test-driven with Ginkgo.  We recommend that any new code be test-driven using Ginkgo.

# Vendoring

Vendoring for this project is done using `dep`.
