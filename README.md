## stembuild: build your stemcells

# Requirements

Needs:

* `ovftool` (VMWare) on your path 
* clone librsync: https://github.com/charlievieth/librsync/tree/mingw64-fseeko64-v2.0.0


# Build Instructions
Note: you can only build on your host OS. you cannot cross compile. 

1) download cmake
1) set source folder to librsync 
1) set destination folder to librsync/build
1) set CMAKE_INSTALL_PREFIX=librsync/install
1) set CMAKE_BUILD_TYPE=release
1) click configure
1) click generate
2) copy all source and header files from librsync/src into rdiff/
2) copy all source and header files from librsync/build/src into rdiff/
3) remove rdiff/rdiff.c


# Instructions

If you have a `.vhd` file and a `.patch` file for that VHD, then:

```
stembuild -vhd my-vhd.vhd -delta patchfile.patch -version 1000.0
```

will create a stemcell with version `1000.0` in your current working direcory.

Process takes between 10 and 20 minutes. See Progress with `-debug` flag.

Other options available with `stembuild -h`.
